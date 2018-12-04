package merge

import (
	"database/sql"
	"log"
)

type Merge struct {
	config *Config   //配置
	master *sql.DB   //主库
	slave  *sql.DB   //从库
	target *sql.DB   //目标库
	tasks  []Tasker  //任务集合
	tp     *TaskPool //任务工作池
}

func New(config *Config) *Merge {
	m := &Merge{}
	m.config = config
	m.tp = NewTaskPool()
	return m
}

func (m *Merge) NewPool() {
	m.tp.ClearTask()
}

//增加一个工作
func (m *Merge) AddWork(w worker) {
	m.tp.AddTask(w)
}

func (m *Merge) Main() error {
	var err error
	for _, v := range m.config.DBInfos {
		switch v.Name {
		case "master":
			m.master, err = sql.Open("mysql", v.DataSource)
			m.master.SetMaxOpenConns(m.config.ThreadInfos.MaxThreads * 2)
			m.master.SetMaxIdleConns(m.config.ThreadInfos.MaxThreads)
			if err != nil {
				return err
			}
		case "slave":
			m.slave, err = sql.Open("mysql", v.DataSource)
			if err != nil {
				return err
			}
		case "target":
			m.target, err = sql.Open("mysql", v.DataSource)
			if err != nil {
				return err
			}
		}
	}

	m.initTask()
	m.runTask()
	return err
}

// 通过名字获取任务
func (m *Merge) Task(name string) interface{} {
	for _, v := range m.tasks {
		if v.Name() == name {
			return v
		}
	}
	return nil
}

// 初始化任务集合，按照顺序加入，执行顺序即为加入的顺序
func (m *Merge) initTask() {
	m.tasks = append(m.tasks,
		NewQueryTableTask(),      //查询待合并的各表的信息
		NewCreateTableTask(),     //在目标库中创建表结构
		NewSyncTableTask(),       //复制主从库数据到目标库
		NewQueryConflictTask(),   //查询主从库的冲突数据
		NewResolveConflictTask(), //解决冲突的数据(除二进制以外的数据)
		NewResolveBinaryTask(),   //解决冲突的二进制数据
		NewMergeTableTask(),      //合并二进制数据
		NewCleanTask(),           //清理临时数据表
		NewSqlTask(),             //执行额外的sql语句
	)
}

// 按顺序执行任务
func (m *Merge) runTask() {
	context := &Context{m}
	log.Printf("total task: %d\n", len(m.tasks))
	for k, v := range m.tasks {
		log.Printf("task %d (%s) begin\n", k, v.Name())
		v.Prepare(context) // 任务初始化
		for {
			m.NewPool()                       //重量工作线程池
			v.Run(context)                    //运行任务(这个函数里，任务增加工作线程)
			m.tp.Run()                        //启动线程池，执行当前任务添加的所有的工作线程
			m.tp.Wait()                       //等待所有工作线程结束
			if err := v.Error(); err != nil { //检查是否有错误产生
				log.Println("panic,", err)
				return
			}
			v.Complete(context)       //通知任务当前所有线程已经执行完毕
			if !v.Continue(context) { //当前任务是否还要继续，任务可以分段执行
				break
			}
		}

		log.Printf("task %d complete\n", k)
	}
	log.Println("all task complete")
}

// 退出
func (m *Merge) Exit() {
	if m.master != nil {
		m.master.Close()
	}
	if m.slave != nil {
		m.slave.Close()
	}
	if m.target != nil {
		m.target.Close()
	}
}
