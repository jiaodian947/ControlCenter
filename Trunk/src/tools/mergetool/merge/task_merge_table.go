package merge

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"tools/mergetool/gameobj"
	"tools/mergetool/utils"
)

// 合并表
type MergeTableTask struct {
	BaseTask
	stage int
}

func NewMergeTableTask() *MergeTableTask {
	s := &MergeTableTask{}
	s.name = MERGETABLE
	s.stage = 0
	return s
}

func (m *MergeTableTask) Run(context *Context) {
	qt := context.merge.Task(QUERYTABLE).(*QueryTableTask)
	for _, v := range context.merge.config.Tables {
		if v.Mode == M_INSERT {
			if _, has := qt.SlaveTableInfo[v.Name]; has {
				if m.stage == 0 { // 第一阶段合并二进制数据
					if context.merge.config.HasMergeTable(v.Name) {
						w := NewMergeTableWork(context, v.Name, v.Name+SLAVE_SUFFIX, context.merge.target)
						w.owner = m
						context.merge.AddWork(w)
					}
				} else { // 第二阶段合并表
					w := NewCopyTableWork(context, v.Name, v.Name+SLAVE_SUFFIX, v.InsertCols, context.merge.target)
					w.owner = m
					context.merge.AddWork(w)
				}
			}
		}
	}
}

func (m *MergeTableTask) Complete(context *Context) {
	m.stage++
}

func (m *MergeTableTask) Continue(context *Context) bool {
	return m.stage < 2
}

// 从表数据直接复制到主表
type CopyTableWork struct {
	BaseWork
	src    string
	target string
	cols   string
	db     *sql.DB
}

func NewCopyTableWork(context *Context, target, src, cols string, db *sql.DB) *CopyTableWork {
	m := &CopyTableWork{}
	m.context = context
	m.target = target
	m.src = src
	m.cols = cols
	m.db = db
	return m
}

func (m *CopyTableWork) Start() {
	start := time.Now()
	count := TableRows(m.db, m.src)
	var sql string
	if m.cols == "" {
		sql = fmt.Sprintf("INSERT INTO `%s` SELECT * FROM `%s`", m.target, m.src)
	} else {
		sql = fmt.Sprintf("INSERT INTO `%s`(%s) SELECT %s FROM `%s`", m.target, m.cols, m.cols, m.src)
	}

	r, err := m.db.Exec(sql)
	if err != nil {
		m.owner.SetError(err)
		return
	}

	if n, _ := r.RowsAffected(); int(n) != count {
		m.owner.SetError(fmt.Errorf("rows not match"))
		return
	}

	log.Printf("copy table %s to %s complete, count: %d. elapsed time: %.2f seconds\n", m.src, m.target, count, time.Now().Sub(start).Seconds())
}

// 合并表的二进制数据
type MergeTableWork struct {
	BaseWork
	src        string
	target     string
	mergeTable MergeTable
	db         *sql.DB
	buff       []byte
}

func NewMergeTableWork(context *Context, target, src string, db *sql.DB) *MergeTableWork {
	m := &MergeTableWork{}
	m.context = context
	m.src = src
	m.target = target
	m.db = db
	m.buff = make([]byte, 0, gameobj.MAX_DATA_LEN)
	return m
}

func (m *MergeTableWork) Start() {
	for k := range m.context.merge.config.Merge {
		if m.context.merge.config.Merge[k].Name == m.target {
			m.mergeTable = m.context.merge.config.Merge[k]
			break
		}
	}

	// 准备数据
	m.mergeTable.Prepare()

	switch m.mergeTable.Mode {
	case M_INSERT:
		m.InsertMode() //插入模式
	case M_MERGE:
		m.MergeMode() //合并模式
	default:
		m.owner.SetError(fmt.Errorf("unsupport mode %s", m.mergeTable.Mode))
	}
}

// 插入模式，插入模式下，如果配置里定义了某一行的处理规则，则对这一行的数据进行特殊处理
// 如果没有配置，则其它的行插入主表。
// 处理方式：
//    1.如果特殊处理的，则处理完成后，在从表里删除这一条数据。
//    2.不需要处理的数据留在从表里。在合并的第二阶段，从表的数据会插入到主表，这样就实现了合并的目的。
func (m *MergeTableWork) InsertMode() {
	for _, c := range m.mergeTable.Columns {
		if c.Key == "" {
			m.owner.SetError(fmt.Errorf("merge table %s column %s key not define", m.mergeTable.Name, c.Name))
			return
		}
		if c.GameData == nil {
			continue
		}
		// 获取主表里对象的二进制数据
		targetbinary, err1 := m.GetBinary(c.Name, m.target, m.mergeTable.KeyName, c.Key)
		if err1 != nil {
			m.owner.SetError(err1)
			return
		}
		// 获取从表里对象的二进制数据
		srcbinary, err2 := m.GetBinary(c.Name, m.src, m.mergeTable.KeyName, c.Key)
		if err2 != nil {
			m.owner.SetError(err2)
			return
		}

		// 合并二进制数据
		binary, err := m.ResolveBinary(targetbinary, srcbinary, c.GameData)
		if err != nil {
			m.owner.SetError(err)
			return
		}

		// 如果主表不为空，更新主表数据，删除从表的数据
		if targetbinary != nil {
			//更新主表的数据
			if err := m.UpdateBinary(c.Name, m.target, m.mergeTable.KeyName, c.Key, binary); err != nil {
				m.owner.SetError(err)
				return
			}
			if srcbinary != nil {
				//清除从表的数据
				if err := m.ClearData(m.src, m.mergeTable.KeyName, c.Key); err != nil {
					m.owner.SetError(err)
					return
				}
			}
		} else if srcbinary != nil { //主表为空，从表不为空，则更新从表，这里不要进行从表的删除，在合并表的时候，会和其它数据一起合并到主表。
			if err := m.UpdateBinary(c.Name, m.src, m.mergeTable.KeyName, c.Key, binary); err != nil {
				m.owner.SetError(err)
				return
			}
		}
	}
}

// 获取某个字段的二进制数据
func (m *MergeTableWork) GetBinary(field, table, keyname, key string) ([]byte, error) {
	sql := fmt.Sprintf("SELECT `%s` FROM `%s` WHERE `%s`='%s'", field, table, keyname, key)
	r, err := m.db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	if !r.Next() {
		return nil, nil
	}
	var binary []byte
	if err = r.Scan(&binary); err != nil {
		return nil, err
	}

	return binary, nil
}

// 更新二进制数据到数据库
func (m *MergeTableWork) UpdateBinary(field, table, keyname, key string, binary []byte) error {
	sql := fmt.Sprintf("UPDATE `%s` SET `%s`=? WHERE `%s`=?", table, field, keyname)
	stmt, err := m.db.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	r, err := stmt.Exec(binary, key)
	if err != nil {
		return err
	}

	if n, _ := r.RowsAffected(); n == 0 {
		return fmt.Errorf("update binary failed")
	}

	return nil
}

//清除数据(一行)
func (m *MergeTableWork) ClearData(table, keyname, key string) error {
	sql := fmt.Sprintf("DELETE FROM `%s` WHERE `%s`=?", table, keyname)
	stmt, err := m.db.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	r, err := stmt.Exec(key)
	if err != nil {
		return err
	}

	if n, _ := r.RowsAffected(); n == 0 {
		return fmt.Errorf("clear binary failed")
	}

	return nil
}

func (m *MergeTableWork) ResolveBinary(target, src []byte, gamedata *GameData) ([]byte, error) {
	// 预处理，将数据的信息写入map，提高查询的效率。
	gamedata.Prepare()
	// 从二进制数据，实例化对象
	targetObj := gameobj.NewGameDataFromBinary(target)
	srcObj := gameobj.NewGameDataFromBinary(src)
	if targetObj == nil || srcObj == nil {
		return nil, fmt.Errorf("resolve binary failed")
	}

	log.Println("target:", gameobj.OutputXml(targetObj))
	log.Println("src:", gameobj.OutputXml(targetObj))

	// 需要处理的属性
	for _, v := range gamedata.GameAttrs {
		attr := targetObj.Attrs.GetAttr(v.Name)
		attr2 := srcObj.Attrs.GetAttr(v.Name)

		switch v.Mode {
		case M_ADD: //值相加
			if attr == nil && attr2 != nil {
				targetObj.Attrs.AddAttr(attr2)
			} else if attr != nil && attr2 != nil {
				if !attr.Add(attr2) {
					return nil, fmt.Errorf("resolve binary attr add failed")
				}
			}
		case M_CLEAR: //值清空
			if attr != nil {
				attr.Clear()
			}
		}
	}

	// 处理表格
	recs := targetObj.Records.NameList()
	for _, recname := range recs {
		mode := gamedata.Mode
		rec := targetObj.Records.Record(recname)
		rec2 := srcObj.Records.Record(recname)
		r := gamedata.GetRecByName(recname)
		if r != nil && r.Mode != "" {
			mode = r.Mode
		}

		switch mode {
		case M_CLEAR: //清空表格
			rec.Clear()
		case M_INSERT: //从表插入主表
			if rec2 != nil {
				len := rec2.RowCount()
				for i := 0; i < len; i++ {
					rec.AddRowValue(-1, rec2.Row(i))
				}
			}
		}

		//是否需要排序
		if r != nil {
			if r.Key != "" {
				var col int
				if err := utils.ParseStrNumber(r.Key, &col); err != nil {
					return nil, err
				}
				rec.Sort(col, r.Sort)
			}
		}

		//进行行更新，如果超出最大行数，则删除多余的行
		rec.Limit()
	}

	// 序列化对象
	log.Println("resolve target:", gameobj.OutputXml(targetObj))
	ar := utils.NewStoreArchiver(m.buff)
	if err := targetObj.Store(ar); err != nil {
		panic(err)
	}

	if !targetObj.Compress { //不需要压缩
		return ar.Data(), nil
	}

	return gameobj.CompressData(ar.Data()), nil
}

// 合并模式，只处理配置的行，如果没有配置则舍弃。
func (m *MergeTableWork) MergeMode() {
	//解决冲突
	for _, c := range m.mergeTable.Columns {
		if c.Key == "" {
			m.owner.SetError(fmt.Errorf("merge table %s column %s key not define", m.mergeTable.Name, c.Name))
			return
		}
		if c.GameData == nil {
			continue
		}

		targetbinary, err1 := m.GetBinary(c.Name, m.target, m.mergeTable.KeyName, c.Key)
		if err1 != nil {
			m.owner.SetError(err1)
			return
		}
		srcbinary, err2 := m.GetBinary(c.Name, m.src, m.mergeTable.KeyName, c.Key)
		if err2 != nil {
			m.owner.SetError(err2)
			return
		}

		binary, err := m.ResolveBinary(targetbinary, srcbinary, c.GameData)
		if err != nil {
			m.owner.SetError(err)
			return
		}

		// 如果主表不为空，更新主表数据，删除从表的数据
		if targetbinary != nil {
			//更新主表的数据
			if err := m.UpdateBinary(c.Name, m.target, m.mergeTable.KeyName, c.Key, binary); err != nil {
				m.owner.SetError(err)
				return
			}
			if srcbinary != nil {
				//清除从表的数据
				if err := m.ClearData(m.src, m.mergeTable.KeyName, c.Key); err != nil {
					m.owner.SetError(err)
					return
				}
			}
		} else if srcbinary != nil { //主表为空，从表不为空，则更新从表，这里不要进行从表的删除，在合并表的时候，会和其它数据一起合并到主表。
			if err := m.UpdateBinary(c.Name, m.src, m.mergeTable.KeyName, c.Key, binary); err != nil {
				m.owner.SetError(err)
				return
			}
		}
	}

	// 删除未配置的行
	if err := m.ClearUndefine(m.target, m.mergeTable.KeyName); err != nil {
		m.owner.SetError(err)
		return
	}
	if err := m.ClearUndefine(m.src, m.mergeTable.KeyName); err != nil {
		m.owner.SetError(err)
		return
	}
}

// 清理未配置的行数据
func (m *MergeTableWork) ClearUndefine(table, keyname string) error {
	sql := fmt.Sprintf("SELECT `%s` FROM `%s`", keyname, table)
	r, err := m.db.Query(sql)
	if err != nil {
		return err
	}
	del := make([]string, 0, 32)
	for r.Next() {
		var key string
		if err := r.Scan(&key); err != nil {
			r.Close()
			return err
		}

		if m.mergeTable.FindKeyIndex(key) == -1 { //没有找到，插入到待删除列表
			del = append(del, key)
		}
	}
	r.Close()

	for _, d := range del {
		err := m.ClearData(table, keyname, d)
		if err != nil {
			return err
		}
	}

	return nil
}
