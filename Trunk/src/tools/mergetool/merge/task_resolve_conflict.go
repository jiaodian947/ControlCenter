package merge

import (
	"database/sql"
	"fmt"
	"log"
)

type ResolveConflictTask struct {
	BaseTask
	stage int
}

func NewResolveConflictTask() *ResolveConflictTask {
	s := &ResolveConflictTask{}
	s.name = RESOLVECONFLICT
	s.stage = 0
	return s
}

func (s *ResolveConflictTask) Run(context *Context) {
	switch s.stage {
	case 0:
		{
			qc := context.merge.Task(QUERYCONFLICT).(*QueryConflictTask)
			// delete
			for _, conflict := range qc.ConflictDels {
				work := NewResloveConflictDelWork(context, conflict.table+SLAVE_SUFFIX, conflict.field, conflict.conflicts)
				context.merge.AddWork(work)
			}
		}
	case 1:
		{
			for k := range context.merge.config.Resolve {
				resolve := context.merge.config.Resolve[k]
				tbl := resolve.Name + SLAVE_SUFFIX
				work := NewResolveConflictUpdateWork(context, tbl, resolve.Columns)
				context.merge.AddWork(work)
			}
		}
	}
}

func (s *ResolveConflictTask) Complete(context *Context) {
	s.stage++
}

func (s *ResolveConflictTask) Continue(context *Context) bool {
	return s.stage < 2
}

type ResolveConflictDelWork struct {
	BaseWork
	db    *sql.DB
	table string
	field string
	id    []string
}

func NewResloveConflictDelWork(context *Context, tbl string, field string, id []string) *ResolveConflictDelWork {
	s := &ResolveConflictDelWork{}
	s.context = context
	s.table = tbl
	s.field = field
	s.id = id
	s.db = context.merge.target
	return s
}

func (s *ResolveConflictDelWork) Start() {
	s.Delete()
}

func (s *ResolveConflictDelWork) Delete() {
	stmt, err := s.db.Prepare(fmt.Sprintf("DELETE FROM `%s` WHERE `%s` = ?", s.table, s.field))
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	count := 0
	for k := range s.id {
		r, err := stmt.Exec(s.id[k])
		if err != nil {
			panic(err)
		}
		row, _ := r.RowsAffected()
		count += int(row)
	}
	log.Printf("delete %s total %d\n", s.table, count)
}

type ResolveConflictUpdateWork struct {
	BaseWork
	db    *sql.DB
	table string
	cols  []ResolveTableCol
}

func NewResolveConflictUpdateWork(context *Context, tbl string, cols []ResolveTableCol) *ResolveConflictUpdateWork {
	s := &ResolveConflictUpdateWork{}
	s.context = context
	s.table = tbl
	s.cols = cols
	s.db = context.merge.target
	return s
}

func (s *ResolveConflictUpdateWork) Start() {
	for k1 := range s.cols {
		column := s.cols[k1]
		if column.GameObj == nil && column.GameData == nil {
			s.Update(column)
		}
	}
}

func (s *ResolveConflictUpdateWork) Update(column ResolveTableCol) {
	qc := s.context.merge.Task(QUERYCONFLICT).(*QueryConflictTask)

	id := column.Value
	field := column.Name
	delete := false
	if id == "" {
		delete = true
		id = qc.TableConflictId(s.table, field)
	}

	log.Println("begin update", s.table, field)
	resolve := qc.ResolveMap(id)
	stmt, err := s.db.Prepare(fmt.Sprintf("UPDATE `%s` SET `%s`=? WHERE `%s`=?", s.table, field, field))
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	count := 0
	for _, v := range resolve {
		var newvalue string
		if delete {
			newvalue = ""
		} else {
			newvalue = v.new
		}
		r, err := stmt.Exec(newvalue, v.old)
		if err != nil {
			panic(err)
		}
		row, _ := r.RowsAffected()
		count += int(row)
	}
	log.Printf("update %s.%s total %d\n", s.table, field, count)
}
