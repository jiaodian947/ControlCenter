package models

import "github.com/astaxie/beego/orm"

func init() {
	orm.RegisterModelWithPrefix("cc_", new(Account))
	orm.RegisterModelWithPrefix("cc_", new(Server))
}
