package models

import (
	"manage/util"

	"github.com/astaxie/beego/orm"
)

type Server struct {
	Id         int
	GameId     int
	ServerName string `orm:"size(128)"`
	ServerId   int
	DistrictId int
	ServerIp   string `orm:"size(32)"`
	ToolPort   int
	GameDb     string `orm:"size(256)"`
	LogDb      string `orm:"size(256)"`
}

func (m *Server) TableName() string {
	return "ftp_common_servers"
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
	return orm.NewOrm().QueryTable("ftp_common_servers").OrderBy("Id")
}
