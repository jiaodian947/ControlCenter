package controllers

import (
	"encoding/json"
	"fmt"
	"speedy/models"
	"strconv"
	"time"
)

type User struct {
	CheckRight
}

func (u *User) AddWhiteListItem() {
	var whiteListData map[string]interface{}
	var reply Reply
	reply.Status = 500

	if err := json.Unmarshal(u.Ctx.Input.RequestBody, &whiteListData); err != nil {
		reply.Data = err.Error()
		u.Data["json"] = &reply
		u.ServeJSON()
		return
	}
	whiteListItem, _ := whiteListData["whiteListItem"].(map[string]interface{})
	data, _ := whiteListItem["data"].(map[string]interface{})
	addTime, _ := strconv.ParseInt(strconv.FormatFloat(whiteListItem["add_time"].(float64), 'f', -1, 64), 10, 64)
	adder := whiteListItem["adder"].(string)
	clentIp := data["ip"].(string)
	note := data["note"].(string)
	whiteList := models.Whitelist{
		ClientIp: clentIp,
		Note:     note,
		AddTime:  time.Unix(addTime/1000, 0),
		Adder:    adder,
	}
	_, err := models.AddWhitelist(&whiteList)
	if err != nil {
		reply.Data = err.Error()
		u.Data["json"] = &reply
		u.ServeJSON()
		return
	}
	reply.Status = 200
	u.Data["json"] = &reply
	u.ServeJSON()
	return
}
func (u *User) DeleteWhiteListItem() {
	var reply Reply
	reply.Status = 500
	id, _ := u.GetInt("id", 0)
	err := models.DeleteWhitelist(id)
	if err != nil {
		reply.Data = err.Error()
		u.Data["json"] = &reply
		u.ServeJSON()
		return
	}
	reply.Status = 200
	reply.Data = map[string]interface{}{}
	u.Data["json"] = &reply
	u.ServeJSON()
}
func (u *User) GetWhiteList() {
	var reply Reply
	reply.Status = 500
	list, err := models.GetAllWhitelist(nil, nil, nil, nil, 0, 0)
	if err != nil {
		reply.Data = err.Error()
		u.Data["json"] = &reply
		u.ServeJSON()
		return
	}
	reply.Status = 200
	reply.Data = map[string]interface{}{
		"whiteList": list,
	}
	u.Data["json"] = &reply
	u.ServeJSON()
}

// get all users's info except password
func (u *User) GetAllUsersInfo() {
	var reply Reply
	reply.Status = 500
	users, err := models.GetAllUser(nil, nil, nil, nil, 0, 0)
	if err != nil {
		reply.Data = err.Error()
		u.Data["json"] = &reply
		u.ServeJSON()
		return
	}
	reply.Status = 200
	reply.Data = users
	u.Data["json"] = &reply
	u.ServeJSON()
}

// get user's info
func (u *User) GetUserInfo() {
	var reply Reply
	reply.Status = 500
	username := u.GetString("token")
	views, err := models.GetUserViewsPower(username)
	if err != nil {
		reply.Data = err.Error()
		u.Data["json"] = &reply
		u.ServeJSON()
		return
	}
	reply.Status = 200
	reply.Data = map[string]interface{}{
		"views":  views,
		"name":   username,
		"avatar": "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
	}
	u.Data["json"] = &reply
	u.ServeJSON()

}

// get user's views power
func (u *User) GetViewsPower() {
	var reply Reply
	reply.Status = 500
	username := u.GetString("username")
	views, err := models.GetUserViewsPower(username)
	if err != nil {
		reply.Data = err.Error()
		u.Data["json"] = &reply
		u.ServeJSON()
		return
	}
	reply.Status = 200
	reply.Data = map[string]interface{}{
		"views": views,
	}
	u.Data["json"] = &reply
	u.ServeJSON()

}

// get user's servers power
func (u *User) GetServersPower() {
	var reply Reply
	reply.Status = 500
	username := u.GetString("username")
	servers, err := models.GetUserServersPower(username)
	if err != nil {
		reply.Data = err.Error()
		u.Data["json"] = &reply
		u.ServeJSON()
		return
	}
	reply.Status = 200
	reply.Data = map[string]interface{}{
		"servers": servers,
	}
	u.Data["json"] = &reply
	u.ServeJSON()

}

// change user's server power
func (u *User) ChangeServerPower() {
	var serverPowerData map[string]interface{}
	var reply Reply
	reply.Status = 500

	if err := json.Unmarshal(u.Ctx.Input.RequestBody, &serverPowerData); err != nil {
		reply.Data = err.Error()
		u.Data["json"] = &reply
		u.ServeJSON()
		return
	}
	serverPower := serverPowerData["server_power"].(map[string]interface{})
	username := serverPower["username"].(string)
	action := serverPower["action"].(string)
	serverPowerArr := serverPower["server_power"].([]interface{})
	fmt.Println(serverPowerArr)
	fmt.Println(username)
	fmt.Println(action)

	switch action {
	case "add":
		for i := 0; i < len(serverPowerArr); i++ {
			item := &models.UsersServersPermission{Username: username, ServerId: int(serverPowerArr[i].(float64))}
			if _, err := models.AddUsersServersPermission(item); err != nil {
				reply.Data = err.Error()
				u.Data["json"] = &reply
				u.ServeJSON()
				return
			}
		}
	case "delete":
		for j := 0; j < len(serverPowerArr); j++ {
			item := &models.UsersServersPermission{Username: username, ServerId: int(serverPowerArr[j].(float64))}
			if _, err := models.DeleteUsersServersPermission(item); err != nil {
				reply.Data = err.Error()
				u.Data["json"] = &reply
				u.ServeJSON()
				return
			}
		}
	default:
		return
	}
	reply.Status = 200
	reply.Data = map[string]interface{}{}
	u.Data["json"] = &reply
	u.ServeJSON()
}

// change user's view power
func (u *User) ChangeViewPower() {
	var viewPowerData map[string]interface{}
	var reply Reply
	reply.Status = 500

	if err := json.Unmarshal(u.Ctx.Input.RequestBody, &viewPowerData); err != nil {
		reply.Data = err.Error()
		u.Data["json"] = &reply
		u.ServeJSON()
		return
	}
	viewPower := viewPowerData["view_power"].(map[string]interface{})
	username := viewPower["username"].(string)
	action := viewPower["action"].(string)
	viewPowerArr := viewPower["view_power"].([]interface{})

	switch action {
	case "add":
		for i := 0; i < len(viewPowerArr); i++ {
			item := &models.UsersViewsPermission{Username: username, ViewId: int(viewPowerArr[i].(float64))}
			if _, err := models.AddUsersViewsPermission(item); err != nil {
				reply.Data = err.Error()
				u.Data["json"] = &reply
				u.ServeJSON()
				return
			}
		}
	case "delete":
		for j := 0; j < len(viewPowerArr); j++ {
			item := &models.UsersViewsPermission{Username: username, ViewId: int(viewPowerArr[j].(float64))}
			if _, err := models.DeleteUsersViewsPermission(item); err != nil {
				reply.Data = err.Error()
				u.Data["json"] = &reply
				u.ServeJSON()
				return
			}
		}
	default:
		return
	}
	reply.Status = 200
	reply.Data = map[string]interface{}{}
	u.Data["json"] = &reply
	u.ServeJSON()

}
