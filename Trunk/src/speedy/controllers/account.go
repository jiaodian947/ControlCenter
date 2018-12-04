package controllers

import (
	"encoding/json"
	"encoding/xml"
	"speedy/models"
	"strconv"
	"time"
)

type Account struct {
	CheckRight
}

func (a *Account) GetSendAccountList() {
	page, _ := a.GetInt64("page", 0)
	pageSize, _ := a.GetInt64("pageSize", 10)
	var count int64
	var err error
	var reply Reply
	reply.Status = 500
	// if page == 0 {
	count, err = models.GetCountSendAccountLog()
	if err != nil {
		reply.Data = err.Error()
		a.Data["json"] = &reply
		a.ServeJSON()
		return
	}
	// }

	mailList, err := models.GetAllSendAccountLog(nil, nil, nil, nil, page*pageSize, pageSize)
	if err != nil {
		reply.Data = err.Error()
		a.Data["json"] = &reply
		a.ServeJSON()
		return
	}
	reply.Status = 200
	reply.Data = map[string]interface{}{
		"accountList": mailList,
		"count":       count,
	}
	a.Data["json"] = &reply
	a.ServeJSON()
}
func (a *Account) DeleteSendAccountLog() {
	var reply Reply
	reply.Status = 500
	id, _ := a.GetInt("id", 0)
	err := models.DeleteSendAccountLog(id)
	if err != nil {
		reply.Data = err.Error()
		a.Data["json"] = &reply
		a.ServeJSON()
		return
	}
	reply.Status = 200
	reply.Data = map[string]interface{}{}
	a.Data["json"] = &reply
	a.ServeJSON()
}
func (a *Account) Ban() {
	var banDetailsData map[string]interface{}
	var reply Reply
	reply.Status = 500

	if err := json.Unmarshal(a.Ctx.Input.RequestBody, &banDetailsData); err != nil {
		reply.Data = err.Error()
		a.Data["json"] = &reply
		a.ServeJSON()
		return
	}
	banDetails := banDetailsData["banDetails"].(map[string]interface{})
	data := banDetails["data"].(map[string]interface{})
	sendReason := data["send_reason"].(string)
	sendType := int(data["send_type"].(float64))
	roles := data["mail_roles"].(string)
	// rolesArr := strings.Split(roles, ",")
	var endTime int64
	endTime = -28800000 //1970-1-1 0:0:0
	if data["end_time"] != "" {
		endTime, _ = strconv.ParseInt(strconv.FormatFloat(data["end_time"].(float64), 'f', -1, 64), 10, 64)

	}
	t := endTime / 1000
	year := time.Unix(t, 0).Year()
	month := int(time.Unix(t, 0).Month())
	day := time.Unix(t, 0).Day()
	hour := time.Unix(t, 0).Hour()
	minute := time.Unix(t, 0).Minute()
	second := time.Unix(t, 0).Second()

	type Block struct {
		XMLName xml.Name `xml:"block"`
		Type    int      `xml:"type,attr"`
		Role    string   `xml:"role,attr"`
		Year    int      `xml:"year,attr"`
		Month   int      `xml:"month,attr"`
		Day     int      `xml:"day,attr"`
		Hour    int      `xml:"hour,attr"`
		Minute  int      `xml:"minute,attr"`
		Second  int      `xml:"second,attr"`
	}
	type CustomXML struct {
		XMLName xml.Name `xml:"custom"`
		Block   []Block  `xml:"block"`
	}

	block := Block{Type: sendType, Role: roles, Year: year, Month: month, Day: day, Hour: hour, Minute: minute, Second: second}
	customxml := new(CustomXML)
	customxml.Block = append(customxml.Block, block)
	customXMLData, _ := xml.Marshal(&customxml)

	sendTime, _ := strconv.ParseInt(strconv.FormatFloat(banDetails["send_time"].(float64), 'f', -1, 64), 10, 64)
	sender := banDetails["sender"].(string)
	gameId := int(banDetails["game_id"].(float64))
	serverId, _ := strconv.Atoi(banDetails["server_id"].(string))
	custom := Custom{GameId: gameId, ServerId: serverId, Type: 25, Custom: string(customXMLData)}
	reply = SendMessageToGameServer(&custom)

	// type SendAccountLog struct {
	// 	Id           int       `orm:"column(id);auto" json:"id"`
	// 	ToRoles      string    `orm:"column(to_roles)" json:"to_roles"`
	// 	Type         int       `orm:"column(type)" json:"type"`
	// 	SendServerid int       `orm:"column(send_serverid);null" json:"send_serverid"`
	// 	SendGameid   int       `orm:"column(send_gameid);null" json:"send_gameid"`
	// 	EndTime      time.Time `orm:"column(end_time);type(datetime);null" json:"end_time"`
	// 	SendTime     time.Time `orm:"column(send_time);type(datetime)" json:"send_time"`
	// 	SendReason   string    `orm:"column(send_reason)" json:"send_reason"`
	// 	Sender       string    `orm:"column(sender)" json:"sender"`
	// }

	sendAccountLog := models.SendAccountLog{
		ToRoles:      roles,
		Type:         sendType,
		SendServerid: serverId,
		SendGameid:   gameId,
		EndTime:      time.Unix(t, 0),
		SendTime:     time.Unix(sendTime/1000, 0),
		SendReason:   sendReason,
		Sender:       sender,
	}
	_, err := models.AddSendAccountLog(&sendAccountLog)
	if err != nil {
		reply.Data = err.Error()
		a.Data["json"] = &reply
		a.ServeJSON()
		return
	}
	reply.Status = 200
	a.Data["json"] = &reply
	a.ServeJSON()
	return
}
