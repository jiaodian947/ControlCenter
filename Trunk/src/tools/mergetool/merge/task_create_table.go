package merge

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

const (
	SLAVE_SUFFIX = "_tmp"
)

// 在目标库中创建对应的表格
type CreateTableTask struct {
	BaseTask
}

func NewCreateTableTask() *CreateTableTask {
	t := &CreateTableTask{}
	t.name = CREATETABLE
	return t
}

// 创建表
func createTable(db *sql.DB, sql, old, tbl string) {
	if old != tbl {
		sql = strings.Replace(sql, old, tbl, -1)
	}
	db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS `%s`", tbl))
	_, err := db.Exec(sql)
	if err != nil {
		log.Fatalf("create %s error, %s", tbl, err)
	}
	log.Printf("create %s success", tbl)
}

func (c *CreateTableTask) Run(context *Context) {
	querytask := context.merge.Task(QUERYTABLE).(*QueryTableTask)
	for _, v := range querytask.MasterTableInfo {
		//创建主库对应的表
		createTable(context.merge.target, v.Create, v.Table, v.Table)
		//创建从库到临时表
		createTable(context.merge.target, v.Create, v.Table, v.Table+SLAVE_SUFFIX)
	}
}

func (c *CreateTableTask) Complete(context *Context) {
	log.Println("all table created")
}
