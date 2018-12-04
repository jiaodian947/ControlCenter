package controllers

import (
	"encoding/json"
	"log"
	"manage/models/gameobj"
	"manage/server"
)

type DatabaseController struct {
	BaseRouter
}

type Response struct {
	Status int
	Err    string
	Reply  *gameobj.GameObject
}

type PlayerRoles struct {
	RoleId   int64
	RoleName string
	Passport string
	SaveData []byte
}

type Query struct {
	GameId   int    `json:"gameid"`
	ServerId int    `json:"serverid"`
	RoleName string `json:"rolename"`
}

func (c *DatabaseController) Post() {
	resp := &Response{}
	resp.Status = 500

	defer func() {
		if p := recover(); p != nil {
			switch inst := p.(type) {
			case error:
				resp.Err = inst.Error()
			case string:
				resp.Err = inst
			}
			resp.Status = 500
			c.Data["json"] = resp
			c.ServeJSON()
		}
	}()

	log.Println(string(c.Ctx.Input.RequestBody))

	var err error
	var query Query
	if err = json.Unmarshal(c.Ctx.Input.RequestBody, &query); err != nil {
		resp.Err = err.Error()
		c.Data["json"] = &resp
		c.ServeJSON()
		log.Println(err)
		return
	}

	log.Println(query)

	gs := server.FindServerByServerId(query.GameId, query.ServerId)
	if gs == nil || gs.DB == nil {
		resp.Err = "server not found"
		c.Data["json"] = resp
		c.ServeJSON()
		return
	}

	row, err := gs.DB.Query("select r.n_roleid, r.s_rolename, r.s_passport, b.lb_save_data from player_roles as r, player_binary as b where r.n_roleid=b.n_roleid and r.s_rolename = ?", query.RoleName)

	if err != nil {
		resp.Err = err.Error()
		c.Data["json"] = resp
		c.ServeJSON()
		return
	}

	defer row.Close()

	r := &PlayerRoles{}
	if row.Next() {
		err := row.Scan(&r.RoleId, &r.RoleName, &r.Passport, &r.SaveData)
		if err != nil {
			resp.Err = err.Error()
			c.Data["json"] = resp
			c.ServeJSON()
			return
		}
	}

	obj := gameobj.NewGameObjectFromBinary(r.SaveData)
	resp.Status = 200
	resp.Reply = obj
	c.Data["json"] = resp
	c.ServeJSON()
}
