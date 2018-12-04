package controllers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"speedy/models"
	"strconv"
	"strings"
	"time"
)

type Mail struct {
	CheckRight
}

func (m *Mail) GetSendMailList() {
	page, _ := m.GetInt64("page", 0)
	pageSize, _ := m.GetInt64("pageSize", 10)
	var count int64
	var err error
	var reply Reply
	reply.Status = 500
	// if page == 0 {
	count, err = models.GetCountSendMailLog()
	if err != nil {
		reply.Data = err.Error()
		m.Data["json"] = &reply
		m.ServeJSON()
		return
	}
	// }

	mailList, err := models.GetAllSendMailLog(nil, nil, nil, nil, page*pageSize, pageSize)
	if err != nil {
		reply.Data = err.Error()
		m.Data["json"] = &reply
		m.ServeJSON()
		return
	}
	reply.Status = 200
	reply.Data = map[string]interface{}{
		"mailList": mailList,
		"count":    count,
	}
	m.Data["json"] = &reply
	m.ServeJSON()
}
func (m *Mail) DeleteSendMailLog() {
	var reply Reply
	reply.Status = 500
	id, _ := m.GetInt("id", 0)
	err := models.DeleteSendMailLog(id)
	if err != nil {
		reply.Data = err.Error()
		m.Data["json"] = &reply
		m.ServeJSON()
		return
	}
	reply.Status = 200
	reply.Data = map[string]interface{}{}
	m.Data["json"] = &reply
	m.ServeJSON()
}
func (m *Mail) SendMail() {
	var mailData map[string]interface{}
	var reply Reply
	reply.Status = 500

	if err := json.Unmarshal(m.Ctx.Input.RequestBody, &mailData); err != nil {
		reply.Data = err.Error()
		m.Data["json"] = &reply
		m.ServeJSON()
		return
	}

	maildetails, _ := mailData["maildetails"].(map[string]interface{})
	data, _ := maildetails["data"].(map[string]interface{})
	attachments, _ := data["attachments"].([]interface{})

	type Obj struct {
		XMLName xml.Name `xml:"obj"`
		Cfg     int      `xml:"cfg,attr"`
		Amount  int      `xml:"count,attr"`
		Bind    int      `xml:"bind,attr"`
	}
	type Mail struct {
		XMLName xml.Name `xml:"mail"`
		Obj     []Obj    `xml:"obj"`
	}
	type Appendix struct {
		XMLName xml.Name `xml:"appendix"`
		Mail    []Mail   `xml:"mail"`
	}
	type CustomXML struct {
		XMLName  xml.Name   `xml:"custom"`
		Type     int        `xml:"type"`
		SubType  int        `xml:"subtype"`
		Title    string     `xml:"title"`
		Content  string     `xml:"content"`
		To       []string   `xml:"to"`
		Appendix []Appendix `xml:"appendix"`
	}

	appendix := Appendix{}
	mail := Mail{}
	for i := 0; i < len(attachments); i++ {
		v, _ := attachments[i].(map[string]interface{})
		cfg := int(v["id"].(float64))
		amount, _ := strconv.Atoi(v["count"].(string))
		bind := v["is_bind"].(bool)
		var isBind = 0
		if bind {
			isBind = 1
		}
		mail.Obj = append(mail.Obj, Obj{Cfg: cfg, Amount: amount, Bind: isBind})
	}

	sendType := int(data["send_type"].(float64))
	sendReason := data["send_reason"].(string)
	mailType, _ := strconv.Atoi(data["mail_type"].(string))
	mailSubType, _ := strconv.Atoi(data["mail_subtype"].(string))
	mailTitle := data["mail_title"].(string)
	mailContent := data["mail_desc"].(string)
	to := data["mail_roles"].(string)
	toArr := strings.Split(to, ",")
	customxml := CustomXML{Type: mailType, SubType: mailSubType, Title: mailTitle, Content: mailContent, To: toArr}
	appendix.Mail = append(appendix.Mail, mail)
	customxml.Appendix = append(customxml.Appendix, appendix)
	customXMLData, _ := xml.Marshal(&customxml)

	gameId := int(maildetails["game_id"].(float64))
	serverId, _ := strconv.Atoi(maildetails["server_id"].(string))
	sendTime, _ := strconv.ParseInt(strconv.FormatFloat(maildetails["send_time"].(float64), 'f', -1, 64), 10, 64)
	sender := maildetails["sender"].(string)
	fmt.Println(sendTime)
	custom := Custom{GameId: gameId, ServerId: serverId, Type: 22, Custom: string(customXMLData)}
	reply = SendMessageToGameServer(&custom)
	// type SendMailLog struct {
	// 	Id           int       `orm:"column(id);auto" json:"id"`
	// 	SendGameid   int       `orm:"column(send_gameid)" json:"send_gameid"`
	// 	SendServerid int       `orm:"column(send_serverid)" json:"send_serverid"`
	// 	ToRoles      string    `orm:"column(to_roles);null" json:"to_roles"`
	// 	MailType     int       `orm:"column(mail_type)" json:"mail_type"`
	// 	MailSubtype  int       `orm:"column(mail_subtype)" json:"mail_subtype"`
	// 	MailTitle    string    `orm:"column(mail_title)" json:"mail_title"`
	// 	MailAppendix string    `orm:"column(mail_appendix);null" json:"mail_appendix"`
	// 	SendType     int       `orm:"column(send_type)" json:"send_type"`
	// 	SendReason   string    `orm:"column(send_reason)" json:"send_reason"`
	// 	SendTime     time.Time `orm:"column(send_time);type(datetime)" json:"send_time"`
	// 	Sender       string    `orm:"column(sender)" json:"sender"`
	// }
	sendMailLog := models.SendMailLog{
		SendGameid:   gameId,
		SendServerid: serverId,
		ToRoles:      to,
		MailType:     mailType,
		MailSubtype:  mailSubType,
		MailTitle:    mailTitle,
		MailAppendix: string(customXMLData),
		SendType:     sendType,
		SendReason:   sendReason,
		SendTime:     time.Unix(sendTime/1000, 0),
		Sender:       sender,
	}
	_, err := models.AddSendMailLog(&sendMailLog)
	if err != nil {
		reply.Data = err.Error()
		m.Data["json"] = &reply
		m.ServeJSON()
		return
	}
	reply.Status = 200
	m.Data["json"] = &reply
	m.ServeJSON()
	return
}
