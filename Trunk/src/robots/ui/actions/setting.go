package actions

import (
	"robots/controller"
	"robots/utils"

	"github.com/lunny/tango"
)

type Setting struct {
	RenderBase
	tango.Json
}

func (i *Setting) PostServer() interface{} {
	type jsondata struct {
		ServerIp   string `json:"serverip"`
		ServerPort string `json:"serverport"`
		ServerId   string `json:"serverid"`
	}
	var data jsondata
	err := i.DecodeJSON(&data)
	if err != nil {
		return map[string]interface{}{
			"err:":   err.Error(),
			"status": 500,
		}
	}

	if data.ServerIp == "" || data.ServerPort == "" || data.ServerId == "" {
		return map[string]interface{}{
			"err:":   "args error",
			"status": 500,
		}
	}

	var port int
	if err := utils.ParseStrNumber(data.ServerPort, &port); err != nil {
		return map[string]interface{}{
			"err:":   err.Error(),
			"status": 500,
		}
	}

	controller.SetServerInfo(data.ServerIp, port, data.ServerId)

	return map[string]interface{}{
		"status": 200,
	}
}

func (i *Setting) PostRobot() interface{} {
	type jsondata struct {
		AccPrefix  string `json:"account_prefix"`
		AccStart   string `json:"account_start"`
		Password   string `json:"password"`
		NamePrefix string `json:"name_prefix"`
		RobotCount string `json:"robot_count"`
	}
	var data jsondata
	err := i.DecodeJSON(&data)
	if err != nil {
		return map[string]interface{}{
			"err:":   err,
			"status": 500,
		}
	}

	if data.AccPrefix == "" ||
		data.AccStart == "" ||
		data.Password == "" ||
		data.NamePrefix == "" ||
		data.RobotCount == "" {
		return map[string]interface{}{
			"err:":   "args error",
			"status": 500,
		}
	}

	var start, count int
	if err := utils.ParseStrNumber(data.AccStart, &start); err != nil {
		return map[string]interface{}{
			"err:":   err.Error(),
			"status": 500,
		}
	}
	if err := utils.ParseStrNumber(data.RobotCount, &count); err != nil {
		return map[string]interface{}{
			"err:":   err.Error(),
			"status": 500,
		}
	}

	controller.SetAccInfo(data.AccPrefix, data.Password, data.NamePrefix, start, count)

	return map[string]interface{}{
		"status": 200,
	}
}
