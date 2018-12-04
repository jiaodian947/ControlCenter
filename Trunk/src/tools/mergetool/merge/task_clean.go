package merge

import (
	"fmt"
	"log"
)

// 清理临时数据表(从表是当作临时表来处理的)
type CleanTask struct {
	BaseTask
}

func NewCleanTask() *CleanTask {
	s := &CleanTask{}
	s.name = CLEAN
	return s
}

func (c *CleanTask) Run(context *Context) {
	querytask := context.merge.Task(QUERYTABLE).(*QueryTableTask)
	for _, v := range querytask.MasterTableInfo {
		sql := fmt.Sprintf("DROP TABLE IF EXISTS `%s`", v.Table+SLAVE_SUFFIX)
		if _, err := context.merge.target.Exec(sql); err != nil {
			c.SetError(err)
			return
		}
		log.Println("drop table", v.Table+SLAVE_SUFFIX)
	}
}
