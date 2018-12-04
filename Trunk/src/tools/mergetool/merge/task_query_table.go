package merge

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
)

// 查询表信息任务
// 查询主从库的表信息，并进行表结构检查，如果表格构不一致，则报错退出
type QueryTableTask struct {
	sync.Mutex
	BaseTask
	MasterTableInfo map[string]*TableInfo //主表的表信息
	SlaveTableInfo  map[string]*TableInfo //从表的表信息
}

func NewQueryTableTask() *QueryTableTask {
	q := &QueryTableTask{}
	q.name = QUERYTABLE
	q.MasterTableInfo = make(map[string]*TableInfo)
	q.SlaveTableInfo = make(map[string]*TableInfo)
	return q
}

func (q *QueryTableTask) Table(table string) *TableInfo {
	if t, has := q.MasterTableInfo[table]; has {
		return t
	}
	return nil
}

// 将所有的表信息进行汇总(由每个任务进行回调)
func (q *QueryTableTask) AddTableInfo(master bool, tbl *TableInfo) {
	q.Lock()
	defer q.Unlock()
	if master {
		q.MasterTableInfo[tbl.Table] = tbl
	} else {
		q.SlaveTableInfo[tbl.Table] = tbl
	}
}

func (q *QueryTableTask) Run(context *Context) {
	for _, v := range context.merge.config.Tables {
		context.merge.AddWork(NewQueryTableWorker(v.Name, true, context, q)) //只从主库加载建库脚本
		context.merge.AddWork(NewQueryTableWorker(v.Name, false, context, q))
	}
}

func (q *QueryTableTask) Complete(context *Context) { //比较两个库表结构是否相等
	for k, v := range q.MasterTableInfo {
		if st, has := q.SlaveTableInfo[k]; has {
			for i := range v.Fields {
				if !v.Fields[i].Equal(st.Fields[i]) {
					fmt.Println(v.Fields[i])
					fmt.Println(st.Fields[i])
					panic(fmt.Sprintf("master and slave table(%s) has different column ", k))
				}
			}
		}
	}
}

type QueryTableWorker struct {
	BaseWork
	tbl    string  //表名
	master bool    //是否是主表
	db     *sql.DB //数据库
}

func NewQueryTableWorker(tbl string, master bool, context *Context, owner *QueryTableTask) *QueryTableWorker {
	w := &QueryTableWorker{}
	w.tbl = tbl
	w.master = master
	w.owner = owner
	if w.master {
		w.db = context.merge.master
	} else {
		w.db = context.merge.slave
	}
	w.context = context
	return w
}

func (w *QueryTableWorker) Start() {
	var dbname string
	if w.master {
		dbname = "[M]"
	} else {
		dbname = "[S]"
	}
	t, err := NewTableInfo(w.db, w.tbl, w.master)
	if err != nil {
		if strings.Contains(err.Error(), "table not found") {
			log.Println(dbname, w.tbl, "not found")
			return
		} else {
			panic(err)
		}
	}
	w.owner.(*QueryTableTask).AddTableInfo(w.master, t)
	log.Println(dbname, "read", w.tbl, "done")
}
