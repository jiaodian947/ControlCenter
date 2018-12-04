package merge

import "tools/mergetool/utils"

type worker interface {
	Start()
}

type TaskPool struct {
	utils.WaitGroupWrapper
	tasks []worker
}

func NewTaskPool() *TaskPool {
	tp := &TaskPool{}
	tp.tasks = make([]worker, 0, 16)
	return tp
}

func (tp *TaskPool) ClearTask() {
	tp.tasks = tp.tasks[:0]
}

func (tp *TaskPool) AddTask(t worker) {
	tp.tasks = append(tp.tasks, t)
}

func (tp *TaskPool) Run() {
	for _, v := range tp.tasks {
		tp.Wrap(v.Start)
	}
}
