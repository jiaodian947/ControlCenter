package models

import (
	"controlcenter/modules/utils"
	"controlcenter/setting"
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
)

type User struct {
	Id       int
	UserName string    `orm:"size(30);unique"`
	NickName string    `orm:"size(30)"`
	Password string    `orm:"size(128)"`
	Email    string    `orm:"size(80);unique"`
	IsAdmin  bool      `orm:"index"`
	IsActive bool      `orm:"index"`
	IsForbid bool      `orm:"index"`
	Lang     int       `orm:"index"`
	Rands    string    `orm:"size(10)"`
	Created  time.Time `orm:"auto_now_add"`
	Updated  time.Time `orm:"auto_now"`
}

func (u *User) TableEngine() string {
	return "INNODB"
}

func (m *User) Insert() error {
	m.Rands = GetUserSalt()
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *User) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *User) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *User) Delete() error {
	if _, err := orm.NewOrm().Delete(m); err != nil {
		return err
	}
	return nil
}

func (m *User) String() string {
	return utils.ToStr(m.Id)
}

func (m *User) Link() string {
	return fmt.Sprintf("%suser/%s", setting.AppUrl, m.UserName)
}

func Users() orm.QuerySeter {
	return orm.NewOrm().QueryTable("user").OrderBy("-Id")
}

func init() {
	orm.RegisterModel(new(User))
}

func GetUserSalt() string {
	return utils.GetRandomString(10)
}
