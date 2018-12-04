package merge

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"strings"
	"time"
)

type SyncTableTask struct {
	BaseTask
}

func NewSyncTableTask() *SyncTableTask {
	s := &SyncTableTask{}
	s.name = SYNCTABLE
	return s
}

func (s *SyncTableTask) Run(context *Context) {
	qt := context.merge.Task(QUERYTABLE).(*QueryTableTask)
	for _, v := range context.merge.config.Tables {
		if v.Mode == M_EMPTY {
			continue
		} else if v.Mode == M_MASTER {
			if _, has := qt.MasterTableInfo[v.Name]; has {
				context.merge.AddWork(NewSyncTableWorker(context, true, v.Name))
			}
		} else {
			if _, has := qt.MasterTableInfo[v.Name]; has {
				context.merge.AddWork(NewSyncTableWorker(context, true, v.Name))
			}
			if _, has := qt.SlaveTableInfo[v.Name]; has {
				context.merge.AddWork(NewSyncTableWorker(context, false, v.Name))
			}
		}
	}

}

type SyncTableWorker struct {
	BaseWork
	master bool    //是否是主表
	tbl    string  //源表名
	target string  //目录表名
	src    *sql.DB //源数据库
	dest   *sql.DB //目标数据库
}

func NewSyncTableWorker(context *Context, master bool, tbl string) *SyncTableWorker {
	s := &SyncTableWorker{}
	s.context = context
	s.master = master
	s.tbl = tbl
	if master {
		s.target = tbl
		s.src = context.merge.master
	} else {
		s.target = tbl + SLAVE_SUFFIX
		s.src = context.merge.slave
	}
	s.dest = context.merge.target
	return s
}

func (s *SyncTableWorker) Start() {
	start := time.Now()
	count := TableRows(s.src, s.tbl)
	var stmt *sql.Stmt
	loops := int(math.Ceil(float64(count) / 100))
	for i := 0; i < loops; i++ {
		r, err := s.src.Query(fmt.Sprintf("SELECT * FROM `%s` LIMIT ?, ?", s.tbl), i*100, 100)
		if err != nil {
			panic(err)
		}
		cols, _ := r.Columns()
		holder := make([]string, len(cols))
		for i := range cols {
			holder[i] = "?"
		}
		if stmt == nil {
			copysql := strings.Join(
				[]string{fmt.Sprintf("INSERT INTO `%s`(", s.target),
					strings.Join(cols, ","),
					") VALUES (",
					strings.Join(holder, ","),
					")",
				},
				"")
			stmt, err = s.dest.Prepare(copysql)
			if err != nil {
				panic(err)
			}
		}

		for r.Next() {
			values := make([][]byte, len(cols))
			scans := make([]interface{}, len(cols))

			for i := range values {
				scans[i] = &values[i]
			}

			if err := r.Scan(scans...); err != nil {
				panic(err)
			}
			_, err := stmt.Exec(scans...)
			if err != nil {
				panic(err)
			}
		}
		r.Close()
	}
	if stmt != nil {
		stmt.Close()
	}

	log.Printf("copy table %s to %s complete, count: %d. elapsed time: %.2f seconds\n", s.tbl, s.target, count, time.Now().Sub(start).Seconds())
}
