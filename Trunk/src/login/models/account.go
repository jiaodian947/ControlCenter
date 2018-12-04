package models

import (
	"login/util"
	"time"

	"github.com/astaxie/beego/orm"
)

type Account struct {
	Id             int
	LogonId        string    `orm:"size(36)"`
	Account        string    `orm:"size(128);index"`
	From           int       `orm:"default(0)"`
	Password       string    `orm:"size(32)"`
	Sex            byte      `orm:"size(1);null"`
	Birthday       time.Time `orm:"type(date);null"`
	OnlineTime     int       `orm:"null"`
	CreateTime     time.Time `orm:"type(datetime);auto_now_add"`
	LastlogTime    time.Time `orm:"type(datetime);auto_now"`
	LastlogAddr    string    `orm:"size(32);null"`
	ValidTime      time.Time `orm:"type(datetime);auto_now_add"`
	LastlogoutTime time.Time `orm:"type(datetime);auto_now"`
	GmLevel        int       `orm:"default(0)"`
	ChargeMode     string    `orm:"size(32);null"`
	Status         byte      `orm:"default(1)"`
	ServerId       int       `orm:"default(0)"`
	GameId         int       `orm:"default(0)"`
	Points         int       `orm:"default(0)"`
	UsedPoints     int       `orm:"default(0)"`
	LastExec       time.Time `orm:"type(datetime);auto_now"`
	TotalTime      int       `orm:"default(0)"`
}

func (m *Account) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *Account) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Account) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Account) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *Account) String() string {
	return util.ToStr(m.Id)
}

func Accounts() orm.QuerySeter {
	return orm.NewOrm().QueryTable("cc_account").OrderBy("-Id")
}
