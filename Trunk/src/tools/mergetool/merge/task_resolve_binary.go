package merge

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"time"
	"tools/mergetool/gameobj"
	"tools/mergetool/utils"
)

type ResolveBinaryTask struct {
	BaseTask
	resolveIndex int
	pagelimit    int
}

func NewResolveBinaryTask() *ResolveBinaryTask {
	r := &ResolveBinaryTask{}
	r.name = RESOLVEBINARY
	r.resolveIndex = 0
	return r
}

func (r *ResolveBinaryTask) Run(context *Context) {
	if r.resolveIndex < len(context.merge.config.Resolve) {
		resolve := context.merge.config.Resolve[r.resolveIndex]
		count := TableRows(context.merge.target, resolve.Name+SLAVE_SUFFIX)
		if count == 0 {
			return
		}
		r.pagelimit = int(math.Ceil(float64(count) / float64(context.merge.config.ThreadInfos.MaxThreads)))

		for k1 := range resolve.Columns {
			column := resolve.Columns[k1]
			if column.GameObj != nil || column.GameData != nil {
				for i := 0; i < context.merge.config.ThreadInfos.MaxThreads; i++ {
					w := NewResolveBinaryWork(context, r, resolve.Name, column.Name, column.GameObj, column.GameData, i*r.pagelimit, r.pagelimit)
					context.merge.AddWork(w)
				}
			}
		}
	}
}

func (r *ResolveBinaryTask) Complete(context *Context) {
	r.resolveIndex++
}

func (r *ResolveBinaryTask) Continue(context *Context) bool {
	return r.resolveIndex < len(context.merge.config.Resolve)
}

type ResolveBinaryWork struct {
	BaseWork
	table    string
	field    string
	db       *sql.DB
	start    int
	limit    int
	object   *GameObj
	data     *GameData
	updateCh chan updateinfo
	quitCh   chan struct{}
	done     bool
}

type updateinfo struct {
	key []byte
	obj Storer
}

func NewResolveBinaryWork(context *Context, owner *ResolveBinaryTask, table, field string, object *GameObj, data *GameData, start, limit int) *ResolveBinaryWork {
	r := &ResolveBinaryWork{}
	r.context = context
	r.owner = owner
	r.table = table
	r.field = field
	r.object = object
	r.data = data
	r.start = start
	r.limit = limit
	r.db = context.merge.target
	r.updateCh = make(chan updateinfo, 32)
	r.quitCh = make(chan struct{})
	r.done = false

	return r
}

func (r *ResolveBinaryWork) Start() {
	q := r.context.merge.Task(QUERYTABLE).(*QueryTableTask)
	key := q.Table(r.table).OneKey()
	if key == "" {
		panic("not found key")
	}
	_sql := fmt.Sprintf("SELECT `%s`,`%s` FROM `%s` Limit %d,%d", key, r.field, r.table+SLAVE_SUFFIX, r.start, r.limit)
	rows, err := r.db.Query(_sql)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	go r.writeCheck(r.table+SLAVE_SUFFIX, key)
	for rows.Next() {
		values := make([][]byte, 2)
		if err := rows.Scan(&values[0], &values[1]); err != nil {
			r.owner.SetError(err)
			break
		}
		if r.object != nil {
			obj := gameobj.NewGameObjectFromBinary(values[1])
			if obj == nil {
				r.owner.SetError(fmt.Errorf("err key: %s", string(values[0])))
				break
			}
			dirty, err := r.ResolveObject(obj)
			if err != nil {
				r.owner.SetError(err)
				break
			}

			if dirty {
				ui := updateinfo{make([]byte, len(values[0])), obj}
				copy(ui.key, values[0])
				r.updateCh <- ui
				continue
			}

		} else if r.data != nil {
			data := gameobj.NewGameDataFromBinary(values[1])
			if data == nil {
				r.owner.SetError(fmt.Errorf("err key: %s", string(values[0])))
				break
			}
			dirty, err := r.ResolveData(data)
			if err != nil {
				r.owner.SetError(err)
				break
			}

			if dirty {
				ui := updateinfo{make([]byte, 0, len(values[0])), data}
				copy(ui.key, values[0])
				r.updateCh <- ui
				continue
			}
		}
	}
	r.done = true
	<-r.quitCh
}

func (r *ResolveBinaryWork) ResolveObject(obj *gameobj.GameObject) (bool, error) {
	qc := r.context.merge.Task(QUERYCONFLICT).(*QueryConflictTask)
	dirty := false
	for _, v := range r.object.GameAttrs {
		id := v.Value
		old := obj.Attrs.GetAttr(v.Name)
		if old == nil {
			log.Println(v.Name, "not found")
			continue
		}
		if id == "" {
			old.Clear()
			dirty = true
			continue
		}
		oldval := old.ToString()
		newval := qc.Rename(id, oldval)
		if newval == "" {
			continue
		}
		if err := old.FromString(newval); err != nil {
			return false, err
		}
		dirty = true
	}

	for _, r := range r.object.GameRecs {
		rec := obj.Records.Record(r.Name)
		if rec == nil {
			continue
		}
		for _, c := range r.Cols {
			if c.Index < 0 || c.Index >= rec.Cols() {
				continue
			}
			id := c.Value
			rows := rec.RowCount()
			for i := 0; i < rows; i++ {
				row := rec.Row(i)
				old := row.Value(c.Index)
				if old == nil {
					continue
				}
				if id == "" {
					old.Clear()
					dirty = true
					continue
				}
				oldval := old.ToString()
				newval := qc.Rename(id, oldval)
				if newval == "" {
					continue
				}

				if err := old.FromString(newval); err != nil {
					return false, err
				}
				dirty = true
			}
		}
	}
	return dirty, nil
}

func (r *ResolveBinaryWork) ResolveData(data *gameobj.GameData) (bool, error) {
	qc := r.context.merge.Task(QUERYCONFLICT).(*QueryConflictTask)
	dirty := false
	for _, v := range r.object.GameAttrs {
		id := v.Value
		old := data.Attrs.GetAttr(v.Name)
		if old == nil {
			log.Println(v.Name, "not found")
			continue
		}
		oldval := old.ToString()
		newval := qc.Rename(id, oldval)
		if newval != "" {
			if err := old.FromString(newval); err != nil {
				return false, err
			}
			dirty = true
			log.Println("replace", v.Name, oldval, "=>", old.ToString())
		}

	}
	return dirty, nil
}

func (r *ResolveBinaryWork) writeCheck(table, key string) {
	sql := fmt.Sprintf("UPDATE `%s` SET `%s`=? WHERE `%s`=?", table, r.field, key)
	stmt, err := r.db.Prepare(sql)
	if err != nil {
		log.Fatalln(err)
		close(r.quitCh)
		return
	}
	count := 0
	buff := make([]byte, 0, gameobj.MAX_DATA_LEN)
	defer func() {
		close(r.quitCh)
		stmt.Close()
		log.Println("update binary", table, r.field, "total", count)
	}()

	for !r.done {
	loop:
		for {
			select {
			case o := <-r.updateCh:
				r.UpdateBinary(stmt, o.key, o.obj, buff)
				count++
			default:
				break loop
			}
		}
		time.Sleep(time.Millisecond)
	}

	//检查处理队列是否为空
	for {
		select {
		case o := <-r.updateCh:
			r.UpdateBinary(stmt, o.key, o.obj, buff)
			count++
		default:
			return
		}
	}
}

type Storer interface {
	Store(ar *utils.StoreArchive) error
	NeedCompress() bool
}

func (r *ResolveBinaryWork) UpdateBinary(stmt *sql.Stmt, key []byte, obj Storer, buff []byte) {
	ar := utils.NewStoreArchiver(buff)
	if err := obj.Store(ar); err != nil {
		panic(err)
	}
	if obj.NeedCompress() {
		//压缩数据
		b := gameobj.CompressData(ar.Data())

		res, err := stmt.Exec(b, key)
		if err != nil {
			log.Println(err)
			panic(err)
		}
		if n, _ := res.RowsAffected(); n == 0 {
			panic("row affected is 0")
		}
	} else {
		res, err := stmt.Exec(ar.Data(), key)
		if err != nil {
			log.Println(err)
			panic(err)
		}
		if n, _ := res.RowsAffected(); n == 0 {
			panic("row affected is 0")
		}
	}

}
