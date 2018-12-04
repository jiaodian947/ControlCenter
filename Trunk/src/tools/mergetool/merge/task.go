package merge

const (
	QUERYTABLE      = "QueryTable"
	CREATETABLE     = "CreateTable"
	SYNCTABLE       = "SyncTable"
	QUERYCONFLICT   = "QueryConflict"
	RESOLVECONFLICT = "ResolveConflict"
	RESOLVEBINARY   = "ResolveBinary"
	MERGETABLE      = "MergeTable"
	CLEAN           = "Clean"
	SQLEXEC         = "SQL Exec"
)

type Tasker interface {
	Name() string
	Prepare(context *Context)
	Run(context *Context)
	Complete(context *Context)
	Continue(context *Context) bool
	Error() error
	SetError(err error)
}

type BaseTask struct {
	name string
	err  error
}

func (b *BaseTask) Name() string {
	return b.name
}

func (b *BaseTask) Prepare(context *Context) {
}

func (b *BaseTask) Run(context *Context) {
}

func (b *BaseTask) Complete(context *Context) {
}

func (b *BaseTask) Continue(context *Context) bool {
	return false
}

func (b *BaseTask) Error() error {
	return b.err
}

func (b *BaseTask) SetError(err error) {
	b.err = err
}

type BaseWork struct {
	context *Context
	owner   Tasker
}
