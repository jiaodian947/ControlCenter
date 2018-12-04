package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type ActivationCode struct {
	Code    string `orm:"size(32);pk"`
	Used    bool
	Updated time.Time `orm:"auto_now"`
}

func (a *ActivationCode) Insert() error {

	if _, err := orm.NewOrm().Insert(a); err != nil {
		return err
	}
	return nil
}

func (a *ActivationCode) Read(fields ...string) error {
	if err := orm.NewOrm().Read(a, fields...); err != nil {
		return err
	}
	return nil
}

func (a *ActivationCode) Use() error {
	a.Used = true
	if _, err := orm.NewOrm().Update(a, "Used", "Updated"); err != nil {
		return err
	}

	return nil
}

func init() {
	orm.RegisterModel(new(ActivationCode))
}
