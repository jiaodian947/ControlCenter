package merge

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"tools/mergetool/utils"
)

// 冲突的数据
type Conflict struct {
	id        string //冲突的id
	name      string //字段名
	old       string //原始值
	value     string //替换的格式
	rfunc     string //替换的函数
	condition string //需要满足的条件
}

type ConflictResolve struct {
	id  string //冲突的id
	old string //原值
	new string //替换后的值
}

//冲突的结果
type ResolveMap map[string]*ConflictResolve

func (r ResolveMap) Has(id string) bool {
	if _, has := r[id]; has {
		return true
	}
	return false
}

// 待删除的冲突信息
type ConflictDelete struct {
	table     string
	field     string
	conflicts []string
}

type FieldConflictDict map[string]*Conflict         // field -> conflict
type TableConflictDict map[string]FieldConflictDict // table -> field
type TableConflictRefer map[string]string           // 表格引用

//获取表，字段的冲突信息
func (t TableConflictDict) Conflict(table string, field string) *Conflict {
	if fdic, has := t[table]; has {
		if c, has1 := fdic[field]; has1 {
			return c
		}
	}
	return nil
}

// 查询冲突任务
type QueryConflictTask struct {
	sync.Mutex
	BaseTask
	ConflictDict  TableConflictDict //冲突的定义
	TableRefer    TableConflictRefer
	IdConflictMap map[string]ResolveMap // key为id
	ConflictDels  []*ConflictDelete     //需要删除的冲突
}

func NewQueryConflictTask() *QueryConflictTask {
	q := &QueryConflictTask{}
	q.name = QUERYCONFLICT
	q.IdConflictMap = make(map[string]ResolveMap)
	q.ConflictDict = make(TableConflictDict)
	q.ConflictDels = make([]*ConflictDelete, 0, 16)
	q.TableRefer = make(TableConflictRefer)
	return q
}

// 是否有某个冲突的定义
func (q *QueryConflictTask) HasConflict(id string) bool {
	if _, has := q.IdConflictMap[id]; has {
		return true
	}
	return false
}

// 处理冲突的定义
func (q *QueryConflictTask) Prepare(context *Context) {
	for _, v := range context.merge.config.Conflict {
		if v.Refer != "" {
			q.TableRefer[v.Name] = v.Refer
		}
		field := make(FieldConflictDict)
		for _, c := range v.Columns {
			f := &Conflict{}
			f.name = c.Name
			f.id = c.Id
			f.value = c.Value
			f.condition = c.Condition
			f.rfunc = c.Func
			if _, dup := field[c.Name]; dup {
				panic("field duplicate")
			}
			field[c.Name] = f
			if _, dup := q.IdConflictMap[c.Id]; dup {
				panic("id duplicate")
			}
			q.IdConflictMap[c.Id] = make(ResolveMap, 1024)
		}
		q.ConflictDict[v.Name] = field
	}
}

func (q *QueryConflictTask) Run(context *Context) {
	qt := context.merge.Task(QUERYTABLE).(*QueryTableTask)
	for _, v := range context.merge.config.Tables {
		if v.Mode != "insert" {
			continue
		}
		if t, has := qt.MasterTableInfo[v.Name]; has {
			if _, has1 := qt.SlaveTableInfo[v.Name]; has1 { // 主表和从表是否存在
				for _, field := range t.Keys {
					context.merge.AddWork(NewQueryConflictWork(context, q, v.Name, field))
				}
			}
		}
	}
}

func (q *QueryConflictTask) Complete(context *Context) {
	for k, v := range q.IdConflictMap {
		log.Printf("conflict `%s` found, total %d\n", k, len(v))
	}

	for _, v := range q.ConflictDels {
		log.Printf("table %s column %s has delete conflicts, total %d\n", v.table, v.field, len(v.conflicts))
	}

}

// 存储所有的冲突
// 有两种存储方式，如果定义冲突解决方案，则按id=>conflict的形式，如果没有定义，则放入删除列表
func (q *QueryConflictTask) ReceiveConflict(table, field string, conflicts []string) {
	q.Lock()
	defer q.Unlock()
	conflict := q.ConflictDict.Conflict(table, field)
	if conflict == nil { //查找引用是否存在
		if refer, has := q.TableRefer[table]; has {
			conflict = q.ConflictDict.Conflict(refer, field)
		}
	}
	if conflict != nil && q.HasConflict(conflict.id) {
		for k := range conflicts {
			if q.IdConflictMap[conflict.id].Has(conflicts[k]) {
				continue
			}
			cf := &ConflictResolve{}
			cf.id = conflict.id
			cf.old = conflicts[k]
			if conflict.value != "" {
				cf.new = fmt.Sprintf(conflict.value, cf.old)
			} else if conflict.rfunc != "" {
				switch conflict.rfunc {
				case "createuid":
					cf.new = utils.NewRoleId(cf.old)
				default:
					panic("undefined func")
				}
			} else {
				panic(fmt.Sprintf("table %s column %s conflict not define", table, field))
			}
			q.IdConflictMap[conflict.id][cf.old] = cf
		}
	} else {
		cd := &ConflictDelete{}
		cd.table = table
		cd.field = field
		cd.conflicts = make([]string, len(conflicts))
		for k := range conflicts {
			cd.conflicts[k] = conflicts[k]
		}
		q.ConflictDels = append(q.ConflictDels, cd)
	}
}

// 获取重命名的值
func (q *QueryConflictTask) Rename(id string, old string) string {
	if resolvemap, has := q.IdConflictMap[id]; has {
		if resolve, has1 := resolvemap[old]; has1 {
			return resolve.new
		}
	}
	return ""
}

// 通过id获取解决的信息
func (q *QueryConflictTask) ResolveMap(id string) ResolveMap {
	if resolvemap, has := q.IdConflictMap[id]; has {
		return resolvemap
	}
	return nil
}

// 获取某个表和字段的冲突ID
func (q *QueryConflictTask) TableConflictId(table, field string) string {
	c := q.ConflictDict.Conflict(table, field)
	if c != nil {
		return c.id
	}

	if refer, has := q.TableRefer[table]; has {
		c = q.ConflictDict.Conflict(refer, field)
		if c != nil {
			return c.id
		}
	}
	return ""
}

type QueryConflictWork struct {
	BaseWork
	table     string
	field     string
	db        *sql.DB
	conflicts []string
}

func NewQueryConflictWork(context *Context, owner *QueryConflictTask, table string, field string) *QueryConflictWork {
	q := &QueryConflictWork{}
	q.table = table
	q.field = field
	q.context = context
	q.db = context.merge.target
	q.owner = owner
	q.conflicts = make([]string, 0, 1024)
	return q
}

func (q *QueryConflictWork) Start() {
	//通过内联查找所有冲突
	sql := fmt.Sprintf("SELECT `%s`.`%s` FROM `%s` INNER JOIN `%s%s` ON `%s`.`%s` = `%s%s`.`%s`",
		q.table, q.field, q.table, q.table, SLAVE_SUFFIX,
		q.table, q.field, q.table, SLAVE_SUFFIX, q.field)
	r, err := q.db.Query(sql)
	if err != nil {
		panic(err)
	}

	defer r.Close()
	for r.Next() {
		var val []byte
		err := r.Scan(&val)
		if err != nil {
			panic(err)
		}
		q.conflicts = append(q.conflicts, string(val))
	}

	q.owner.(*QueryConflictTask).ReceiveConflict(q.table, q.field, q.conflicts)
}
