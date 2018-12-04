package merge

import "log"

type SqlTask struct {
	BaseTask
}

func NewSqlTask() *SqlTask {
	t := &SqlTask{}
	t.name = SQLEXEC
	return t
}

func (s *SqlTask) Run(context *Context) {
	db := context.merge.target
	for k := range context.merge.config.Sqls {
		sql := context.merge.config.Sqls[k].Sql
		_, err := db.Exec(sql)
		if err != nil {
			s.SetError(err)
			return
		}
		log.Println("exec", sql)
	}
}
