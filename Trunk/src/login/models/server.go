package models

import (
	"login/util"

	"github.com/astaxie/beego/orm"
)

type Server struct {
	Id       int
	ServerId int
	GameId   int
	ServerIp string `orm:"size(32)"`
	Logged   int8   `orm:"default(0)"`
}

func (m *Server) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *Server) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Server) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Server) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *Server) String() string {
	return util.ToStr(m.Id)
}

func Servers() orm.QuerySeter {
	return orm.NewOrm().QueryTable("cc_server").OrderBy("-Id")
}
