package controllers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"speedy/models"
	"speedy/server"
	"strconv"
	"time"
)

type Activity struct {
	CheckRight
}

func (a *Activity) GetNoticeList() {
	var reply Reply
	reply.Status = 500
	page, _ := a.GetInt64("page", 0)
	pageSize, _ := a.GetInt64("pageSize", 10)

	noticeList, err := models.GetAllNotice(nil, nil, nil, nil, page*pageSize, pageSize)
	if err != nil {
		reply.Data = err.Error()
		a.Data["json"] = &reply
		a.ServeJSON()
		return
	}
	reply.Status = 200
	reply.Data = map[string]interface{}{
		"list": noticeList,
	}
	a.Data["json"] = &reply
	a.ServeJSON()
}
func (a *Activity) AddNotice() {
	var noticeData map[string]interface{}
	var reply Reply
	reply.Status = 500

	if err := json.Unmarshal(a.Ctx.Input.RequestBody, &noticeData); err != nil {
		reply.Data = err.Error()
		a.Data["json"] = &reply
		a.ServeJSON()
		return
	}
	noticeDetails := noticeData["noticeDetails"].(map[string]interface{})

	content := noticeDetails["content"].([]interface{})
	content_str_byte, _ := json.Marshal(content)
	content_str := string(content_str_byte)

	version := int(noticeDetails["version"].(float64))
	createTime, _ := strconv.ParseInt(strconv.FormatFloat(noticeDetails["create_time"].(float64), 'f', -1, 64), 10, 64)
	startTime, _ := strconv.ParseInt(strconv.FormatFloat(noticeDetails["start_time"].(float64), 'f', -1, 64), 10, 64)
	endTime, _ := strconv.ParseInt(strconv.FormatFloat(noticeDetails["end_time"].(float64), 'f', -1, 64), 10, 64)
	noticeItem := models.Notice{NoticeContent: content_str, NoticeVersion: version, CreateTime: time.Unix(createTime/1000, 0), NoticeStartTime: time.Unix(startTime/1000, 0), NoticeEndTime: time.Unix(endTime/1000, 0)}
	// \\192.168.1.180\html\1.0.0\UnityEditor\Config.json
	_, err := models.AddNotice(&noticeItem)
	if err != nil {
		reply.Data = err.Error()
		a.Data["json"] = &reply
		a.ServeJSON()
		return
	}

	// fmt.Println(activityItem)

	reply.Status = 200
	a.Data["json"] = &reply
	a.ServeJSON()
	return
}
func (a *Activity) DeleteNotice() {
	var reply Reply
	reply.Status = 500
	id, _ := a.GetInt("id", 0)
	err := models.DeleteNotice(id)
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
func (a *Activity) GetActivityAnnouncementList() {
	var reply Reply
	reply.Status = 500
	page, _ := a.GetInt64("page", 0)
	pageSize, _ := a.GetInt64("pageSize", 10)

	announcementList, err := models.GetAllActivityAnnouncement(nil, nil, nil, nil, page*pageSize, pageSize)
	if err != nil {
		reply.Data = err.Error()
		a.Data["json"] = &reply
		a.ServeJSON()
		return
	}
	reply.Status = 200
	reply.Data = map[string]interface{}{
		"list": announcementList,
	}
	a.Data["json"] = &reply
	a.ServeJSON()
}
func (a *Activity) AddActivityAnnouncement() {
	var acivityAnnouncementData map[string]interface{}
	var reply Reply
	reply.Status = 500

	if err := json.Unmarshal(a.Ctx.Input.RequestBody, &acivityAnnouncementData); err != nil {
		reply.Data = err.Error()
		a.Data["json"] = &reply
		a.ServeJSON()
		return
	}
	announcementDetails := acivityAnnouncementData["announcementDetails"].(map[string]interface{})
	file_list := announcementDetails["file_list"].([]interface{})
	file_list_str_byte, _ := json.Marshal(file_list)
	file_list_str := string(file_list_str_byte)
	version := int(announcementDetails["version"].(float64))
	createTime, _ := strconv.ParseInt(strconv.FormatFloat(announcementDetails["create_time"].(float64), 'f', -1, 64), 10, 64)
	startTime, _ := strconv.ParseInt(strconv.FormatFloat(announcementDetails["start_time"].(float64), 'f', -1, 64), 10, 64)
	endTime, _ := strconv.ParseInt(strconv.FormatFloat(announcementDetails["end_time"].(float64), 'f', -1, 64), 10, 64)
	activityItem := models.ActivityAnnouncement{ActivityImages: file_list_str, ActivityVersion: version, CreateTime: time.Unix(createTime/1000, 0), ActivityStartTime: time.Unix(startTime/1000, 0), ActivityEndTime: time.Unix(endTime/1000, 0)}
	// \\192.168.1.180\html\1.0.0\UnityEditor\Config.json
	_, err := models.AddActivityAnnouncement(&activityItem)
	if err != nil {
		reply.Data = err.Error()
		a.Data["json"] = &reply
		a.ServeJSON()
		return
	}

	// fmt.Println(activityItem)

	reply.Status = 200
	a.Data["json"] = &reply
	a.ServeJSON()
	return
}
func (a *Activity) DeleteActivityAnnouncement() {
	var reply Reply
	reply.Status = 500
	id, _ := a.GetInt("id", 0)
	err := models.DeleteActivityAnnouncement(id)
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
func (a *Activity) GetSwitchList() {
	var reply Reply
	reply.Status = 500
	serverId, _ := a.GetInt("server_id", 0)
	gameId, _ := a.GetInt("game_id", 0)
	page, _ := a.GetInt64("page", 0)
	pageSize, _ := a.GetInt64("pageSize", 10)
	gameServer := server.FindServerByServerId(gameId, serverId)
	if gameServer == nil || gameServer.DB == nil {
		reply.Data = "server not found"
		a.Data["json"] = &reply
		a.ServeJSON()
		return
	}
	// rows, err := gameServer.DB.Query("SELECT tb_config.s_id,tb_config.s_name,tb_config.n_type,tb_config.s_prop,tb_config.n_status,tb_config.n_version,tb_config.n_operid FROM tb_config limit ?,?", page*10, pageSize)
	rows, err := gameServer.DB.Query("SELECT * FROM tb_config where n_type = 1 limit ?,? ", page*pageSize, pageSize)

	if err != nil {
		reply.Data = err.Error()
		a.Data["json"] = &reply
		a.ServeJSON()
		return
	}
	defer rows.Close()
	// var s_id string
	// var n_type int
	// var s_prop string
	// var s_name string
	// var n_status int
	// var n_version int
	// var n_operid int

	//返回所有列
	cols, _ := rows.Columns()
	//这里表示一行所有列的值，用[]byte表示
	vals := make([][]byte, len(cols))
	//这里表示一行填充数据
	scans := make([]interface{}, len(cols))
	//这里scans引用vals，把数据填充到[]byte里
	for k, _ := range vals {
		scans[k] = &vals[k]
	}

	i := 0
	result := make([]map[string]string, 0, 10)
	for rows.Next() {
		//填充数据
		rows.Scan(scans...)
		//每行数据
		item := make(map[string]string)
		//把vals中的数据复制到row中
		for k, v := range vals {
			key := cols[k]
			//这里把[]byte数据转成string
			item[key] = string(v)
		}
		//放入结果集
		result = append(result, item)
		i++
	}
	reply.Status = 200
	reply.Data = map[string]interface{}{
		"count":      len(result),
		"switchList": result,
	}
	a.Data["json"] = &reply
	a.ServeJSON()
	fmt.Println(result)
	return
}

func (a *Activity) ChangeSwitchStatus() {
	var switchDetailsData map[string]interface{}
	var reply Reply
	reply.Status = 500

	if err := json.Unmarshal(a.Ctx.Input.RequestBody, &switchDetailsData); err != nil {
		reply.Data = err.Error()
		a.Data["json"] = &reply
		a.ServeJSON()
		return
	}
	switchDetails := switchDetailsData["switchDetails"].(map[string]interface{})
	data := switchDetails["data"].(map[string]interface{})
	nStatus := data["n_status"].(bool)
	nType, _ := strconv.Atoi(data["n_type"].(string))
	sId, _ := data["s_id"].(string)
	sName := data["s_name"].(string)
	type Property struct {
		XMLName  xml.Name `xml:"Property"`
		Id       string   `xml:"ID,attr"`
		CurState int      `xml:"CurState,attr"`
		Desc     string   `xml:"desc,attr"`
	}

	type Msg struct {
		XMLName  xml.Name   `xml:"Msg"`
		Type     int        `xml:"Type,attr"`
		Property []Property `xml:"Property"`
	}
	type CustomXML struct {
		XMLName xml.Name `xml:"custom"`
		Msg     []Msg    `xml:"Msg"`
	}
	var curState int
	curState = 0
	if nStatus {
		curState = 1
	}
	property := Property{Id: sId, CurState: curState, Desc: sName}
	msg := Msg{Type: nType}
	msg.Property = append(msg.Property, property)
	customxml := CustomXML{}
	customxml.Msg = append(customxml.Msg, msg)
	customXMLData, _ := xml.Marshal(&customxml)

	//send to game server
	gameId := int(switchDetails["game_id"].(float64))
	serverId, _ := strconv.Atoi(switchDetails["server_id"].(string))
	custom := Custom{GameId: gameId, ServerId: serverId, Type: 29, Custom: string(customXMLData)}
	reply = SendMessageToGameServer(&custom)
	fmt.Println(reply)

	reply.Status = 200
	a.Data["json"] = &reply
	a.ServeJSON()
	return
}
func (a *Activity) ChangeQuestion() {
	var data map[string]interface{}
	var reply Reply
	reply.Status = 500

	if err := json.Unmarshal(a.Ctx.Input.RequestBody, &data); err != nil {
		reply.Data = err.Error()
		a.Data["json"] = &reply
		a.ServeJSON()
		return
	}
	dataDetails := data["data"].(map[string]interface{})
	content := dataDetails["content"].(string)
	version := int(dataDetails["version"].(float64))
	createTime, _ := strconv.ParseInt(strconv.FormatFloat(dataDetails["create_time"].(float64), 'f', -1, 64), 10, 64)
	questionItem := models.Question{QuestionVersion: version, CreateTime: time.Unix(createTime/1000, 0), QuestionContent: content}
	// \\192.168.1.180\html\1.0.0\UnityEditor\Config.json
	_, err := models.AddQuestion(&questionItem)
	if err != nil {
		reply.Data = err.Error()
		a.Data["json"] = &reply
		a.ServeJSON()
		return
	}
	type Property struct {
		XMLName xml.Name `xml:"Property"`
		Content string   `xml:"Content,attr"`
	}

	type Msg struct {
		XMLName  xml.Name   `xml:"Msg"`
		Type     int        `xml:"Type,attr"`
		Property []Property `xml:"Property"`
	}
	type CustomXML struct {
		XMLName xml.Name `xml:"custom"`
		Msg     []Msg    `xml:"Msg"`
	}

	property := Property{Content: content}
	msg := Msg{Type: 3}
	msg.Property = append(msg.Property, property)
	customxml := CustomXML{}
	customxml.Msg = append(customxml.Msg, msg)
	customXMLData, _ := xml.Marshal(&customxml)
	//send to game server
	gameId := int(dataDetails["game_id"].(float64))
	serverId, _ := strconv.Atoi(dataDetails["server_id"].(string))
	custom := Custom{GameId: gameId, ServerId: serverId, Type: 29, Custom: string(customXMLData)}
	reply = SendMessageToGameServer(&custom)
	// fmt.Println(reply)

	reply.Status = 200
	a.Data["json"] = &reply
	a.ServeJSON()
	return
}
