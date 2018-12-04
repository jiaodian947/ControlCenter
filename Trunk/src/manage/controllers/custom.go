package controllers

import (
	"encoding/json"
	"log"
	"manage/protocol"
	"manage/server"
	"strconv"
)

type Mt2SrvController struct {
	BaseRouter
}

type Custom struct {
	GameId   int    `json:"gameid"`
	ServerId int    `json:"serverid"`
	Type     int    `json:"msgtype"`
	Custom   string `json:"custom"`
}

type Return struct {
	Status int
	Err    string
	Reply  string
}

const (
	CENTER_SERVER_STATE_UNKNOWN                 = iota // 初始状态
	CENTER_SERVER_STATE_INIT_ENV                       // 初始化环境
	CENTER_SERVER_STATE_LOAD_BLK_DEVS                  // 加载被屏蔽的设备号
	CENTER_SERVER_STATE_LOAD_ROLES                     // 加载角色信息
	CENTER_SERVER_STATE_LOAD_SCENES                    // 加载所有普通场景
	CENTER_SERVER_STATE_CAN_OPEN                       // 开启服务状态，需要手动开启
	CENTER_SERVER_STATE_OPENED                         // 已经开启服务服务器开启完毕
	CENTER_SERVER_STATE_CLOSE_GATE_SERVER              // 关闭网关服务器
	CENTER_SERVER_STATE_CLEAR_PLAYERS                  // 清理所有玩家
	CENTER_SERVER_STATE_CLEAR_SCENES                   // 清理所有场景
	CENTER_SERVER_STATE_CLOSE_SCENE_SERVER             // 关闭场景服务器
	CENTER_SERVER_STATE_CLOSE_SHARE_DATA_SERVER        // 关闭共享数据服务器
	CENTER_SERVER_STATE_CLOSE_DB_PROXY_SERVER          // 关闭数据库代理服务器
	CENTER_SERVER_STATE_CLOSED                         // 已经关闭服务
	CENTER_SERVER_STATE_ERROR                          // 出错
)

func (c *Mt2SrvController) Post() {
	var custom Custom
	var res Return
	res.Status = 500

	log.Println(string(c.Ctx.Input.RequestBody))

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &custom); err != nil {
		res.Err = err.Error()
		c.Data["json"] = &res
		c.ServeJSON()
		log.Println(err)
		return
	}

	log.Println(custom)

	msg := protocol.NewMessage(1024)
	arstore := protocol.NewStoreArchiver(msg.Body)
	if custom.Type < protocol.E_MT2SRV_CENTER_MAX {
		arstore.Write(uint8(protocol.E_MT2SRV_TAR_CENTER_SRV))
	} else if custom.Type < protocol.E_MT2SRV_SCENE_MAX {
		arstore.Write(uint8(protocol.E_MT2SRV_TAR_SCENE_SRV))
	} else if custom.Type < protocol.E_MT2SRV_DB_MAX {
		arstore.Write(uint8(protocol.E_MT2SRV_TAR_DB_PROXY_SRV))
	} else {
		arstore.Write(uint8(protocol.E_MT2SRV_TAR_CENTER_SRV))
	}
	arstore.Write(uint8(custom.Type))
	arstore.Write(uint8(protocol.LOGIC_MESSAGE)) // 逻辑处理
	arstore.Write(uint8(protocol.SEND_TO_SCENE)) // 发送给场景
	args := protocol.NewVarMsg(1)
	args.AddString(custom.Custom)
	arstore.Write(args)

	msg.Body = msg.Body[:arstore.Len()]
	msg, err := server.SendMessage(custom.GameId, custom.ServerId, msg, true)
	if msg != nil {
		defer msg.Free()
	}

	if err != nil {
		res.Err = err.Error()
		c.Data["json"] = &res
		c.ServeJSON()
		return
	}

	res.Status = 200
	reply := protocol.NewLoadArchiver(msg.Body)
	_, err = reply.ReadUInt64() // seq
	if err != nil {
		res.Err = err.Error()
		c.Data["json"] = &res
		c.ServeJSON()
		return
	}
	msgid, err := reply.ReadUInt8() // msg
	if err != nil {
		res.Err = err.Error()
		c.Data["json"] = &res
		c.ServeJSON()
	}
	varmsg, err := reply.ReadVarMsg()
	if err != nil {
		res.Err = err.Error()
		c.Data["json"] = &res
		c.ServeJSON()
	}

	switch msgid {
	case protocol.E_SRV2MT_MSG_MONITOR_STATE:
		status := varmsg.Int32Val(0)
		res.Reply = strconv.Itoa(int(status))
	case protocol.E_SRV2MT_MSG_MONITOR_OPEN:
		status := varmsg.Int32Val(0)
		res.Reply = strconv.Itoa(int(status))
	case protocol.E_SRV2MT_MSG_SCENE_MSG:
		if varmsg.Type(0) != protocol.VTYPE_STRING {
			res.Status = 500
			res.Err = "args type not string"
			break
		}

		res.Reply = varmsg.StringVal(0)
	}

	c.Data["json"] = &res
	c.ServeJSON()
}
