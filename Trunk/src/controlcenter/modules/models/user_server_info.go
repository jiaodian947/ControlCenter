package models

import (
	"errors"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	ERR_ACC_NOT_FOUND  = errors.New("account not set")
	UserInfoCollection = "user_servers"
)

type UserServerInfo struct {
	GameId     int
	DistrictId int
	ServerId   int
	RoleName   string
	RoleLevel  int
}

type UserInfo struct {
	Id          bson.ObjectId `bson:"_id"`
	Account     string
	ServerInfos []UserServerInfo
}

func (user *UserInfo) Insert(c *mgo.Collection) error {
	if user.Id == "" {
		user.Id = bson.NewObjectId()
	}
	if user.Account == "" {
		return ERR_ACC_NOT_FOUND
	}

	if user.ServerInfos == nil {
		user.ServerInfos = make([]UserServerInfo, 0, 8)
	}
	return c.Insert(user)
}

func (user *UserInfo) Read(c *mgo.Collection) error {
	if user.Account == "" {
		return ERR_ACC_NOT_FOUND
	}
	return c.Find(bson.M{"account": user.Account}).One(user)
}

func (user *UserInfo) UpsetServerInfo(c *mgo.Collection, info UserServerInfo) error {
	if user.Account == "" {
		return ERR_ACC_NOT_FOUND
	}

	n, err := c.Find(bson.M{
		"account": user.Account,
		"serverinfos": bson.M{
			"$elemMatch": bson.M{
				"gameid":     info.GameId,
				"districtid": info.DistrictId,
				"serverid":   info.ServerId,
				"rolename":   info.RoleName,
			},
		},
	}).Count()

	//  更新
	if err == nil && n >= 1 {
		return c.Update(
			bson.M{
				"account": user.Account,
				"serverinfos": bson.M{
					"$elemMatch": bson.M{
						"gameid":     info.GameId,
						"districtid": info.DistrictId,
						"serverid":   info.ServerId,
						"rolename":   info.RoleName,
					},
				},
			},
			bson.M{
				"$set": bson.M{
					"serverinfos.$.rolelevel": info.RoleLevel,
				},
			})
	}

	// 新增加
	return c.Update(
		bson.M{
			"account": user.Account,
		},
		bson.M{
			"$push": bson.M{
				"serverinfos": info,
			},
		})
}
