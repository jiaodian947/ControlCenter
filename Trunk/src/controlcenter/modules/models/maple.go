package models

import (
	"controlcenter/modules/utils"
	"controlcenter/setting"
	"fmt"

	"github.com/astaxie/beego/orm"
)

type ServerGame struct {
	Id        int    `orm:"auto"`
	Name      string `orm:"size(128)"`
	Comment   string `orm:"size(256)"`
	SecretKey string `orm:"size(32)"`
}

func (sg *ServerGame) Insert() error {
	sg.SecretKey = utils.GetRandomString(32)
	if _, err := orm.NewOrm().Insert(sg); err != nil {
		return err
	}
	return nil
}

func (sg *ServerGame) Read(fields ...string) error {
	if err := orm.NewOrm().Read(sg, fields...); err != nil {
		return err
	}
	return nil
}

func (sg *ServerGame) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(sg, fields...); err != nil {
		return err
	}
	return nil
}

func (sg *ServerGame) Delete(fields ...string) error {
	if _, err := orm.NewOrm().Delete(sg, fields...); err != nil {
		return err
	}
	return nil
}

func (sg *ServerGame) String() string {
	return utils.ToStr(sg.Id)
}

func (sg *ServerGame) Link() string {
	return fmt.Sprintf("%sgame/%d", setting.AppUrl, sg.Id)
}

func Games() orm.QuerySeter {
	return orm.NewOrm().QueryTable("ServerGame").OrderBy("Id")
}

type ServerDistrict struct {
	Id           int    `orm:"auto"`
	GameId       int    `orm:"index"`
	DistrictName string `orm:"size(128)"`
	Comment      string `orm:"size(256)"`
	Group        int    //分组
}

func (d *ServerDistrict) Insert() error {
	if _, err := orm.NewOrm().Insert(d); err != nil {
		return err
	}
	return nil
}

func (d *ServerDistrict) Read(fields ...string) error {
	if err := orm.NewOrm().Read(d, fields...); err != nil {
		return err
	}
	return nil
}

func (d *ServerDistrict) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(d, fields...); err != nil {
		return err
	}
	return nil
}

func (d *ServerDistrict) Delete(fields ...string) error {
	if _, err := orm.NewOrm().Delete(d, fields...); err != nil {
		return err
	}
	return nil
}

func (d *ServerDistrict) String() string {
	return utils.ToStr(d.Id)
}

func (d *ServerDistrict) Link() string {
	return fmt.Sprintf("%sgame/%d/%d", setting.AppUrl, d.GameId, d.Id)
}

func Districts() orm.QuerySeter {
	return orm.NewOrm().QueryTable("ServerDistrict").OrderBy("Id")
}

type ServerInfo struct {
	Id             int `orm:"auto"`
	DistrictId     int `orm:"index"`
	GameId         int
	ServerName     string `orm:"size(128)"`
	ServerType     int
	ServerStatus   int
	PlayerCount    int
	PlayerMaxCount int
	ServerIp       string `orm:"size(64)"`
	ServerPort     int
	Comment        string `orm:"size(256)"`
}

func (s *ServerInfo) Insert() error {
	if _, err := orm.NewOrm().Insert(s); err != nil {
		return err
	}
	return nil
}

func (s *ServerInfo) Read(fields ...string) error {
	if err := orm.NewOrm().Read(s, fields...); err != nil {
		return err
	}
	return nil
}

func (s *ServerInfo) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(s, fields...); err != nil {
		return err
	}
	return nil
}

func (s *ServerInfo) Delete(fields ...string) error {
	if _, err := orm.NewOrm().Delete(s, fields...); err != nil {
		return err
	}
	return nil
}

func (s *ServerInfo) String() string {
	return utils.ToStr(s.Id)
}

func (s *ServerInfo) Link() string {
	return fmt.Sprintf("/game/%d/%d/%d", s.GameId, s.DistrictId, s.Id)
}

func Servers() orm.QuerySeter {
	return orm.NewOrm().QueryTable("ServerInfo").OrderBy("Id")
}

func init() {
	orm.RegisterModel(new(ServerGame), new(ServerDistrict), new(ServerInfo))
}
