package models

import "github.com/astaxie/beego/orm"

func init() {
	//orm.RegisterModelWithPrefix("m_", new(User))
	orm.RegisterModel(new(Server))
}
