package controllers

import (
	"encoding/json"
	"strconv"
	"time"
)

type Logs struct {
	CheckRight
}

func (l *Logs) GetGameLogs() {
	var queryData map[string]interface{}
	var reply Reply
	reply.Status = 500
	if err := json.Unmarshal(l.Ctx.Input.RequestBody, &queryData); err != nil {
		reply.Data = err.Error()
		l.Data["json"] = &reply
		l.ServeJSON()
		return
	}
	queryInfo := queryData["player_info"].(map[string]interface{})
	gameId := int(queryInfo["game_id"].(float64))
	serverId, _ := strconv.Atoi(queryInfo["server_id"].(string))
	roleName := queryInfo["role_name"].(string)
	startTime, _ := strconv.ParseInt(strconv.FormatFloat(queryInfo["start_time"].(float64), 'f', -1, 64), 10, 64)
	endTime, _ := strconv.ParseInt(strconv.FormatFloat(queryInfo["end_time"].(float64), 'f', -1, 64), 10, 64)
	page := int(queryInfo["page"].(float64))
	pageSize := int(queryInfo["page_size"].(float64))
	recordName := queryInfo["record_name"].(string)

	query := QueryLogs{GameId: gameId, ServerId: serverId, RoleName: roleName, StartTime: time.Unix(startTime/1000, 0).Format("2006-01-02 15:04:05"), EndTime: time.Unix(endTime/1000, 0).Format("2006-01-02 15:04:05"), Page: page, PageSize: pageSize}
	reply = QueryGameLogs(&query, recordName)
	reply.Status = 200
	l.Data["json"] = &reply
	l.ServeJSON()

}
func (l *Logs) GetGuildFunctionLogs() {
	var queryData map[string]interface{}
	var reply Reply
	reply.Status = 500
	if err := json.Unmarshal(l.Ctx.Input.RequestBody, &queryData); err != nil {
		reply.Data = err.Error()
		l.Data["json"] = &reply
		l.ServeJSON()
		return
	}
	queryInfo := queryData["guild_info"].(map[string]interface{})
	gameId := int(queryInfo["game_id"].(float64))
	serverId, _ := strconv.Atoi(queryInfo["server_id"].(string))
	guildName := queryInfo["guild_name"].(string)
	startTime, _ := strconv.ParseInt(strconv.FormatFloat(queryInfo["start_time"].(float64), 'f', -1, 64), 10, 64)
	endTime, _ := strconv.ParseInt(strconv.FormatFloat(queryInfo["end_time"].(float64), 'f', -1, 64), 10, 64)
	page := int(queryInfo["page"].(float64))
	pageSize := int(queryInfo["page_size"].(float64))
	recordName := queryInfo["record_name"].(string)

	query := GuildFunctionLogs{GameId: gameId, ServerId: serverId, GuildName: guildName, StartTime: time.Unix(startTime/1000, 0).Format("2006-01-02 15:04:05"), EndTime: time.Unix(endTime/1000, 0).Format("2006-01-02 15:04:05"), Page: page, PageSize: pageSize}
	reply = QueryGuildFunctionLogs(&query, recordName)
	reply.Status = 200
	l.Data["json"] = &reply
	l.ServeJSON()

}
func (l *Logs) GetGuildGameLogs() {
	var queryData map[string]interface{}
	var reply Reply
	reply.Status = 500
	if err := json.Unmarshal(l.Ctx.Input.RequestBody, &queryData); err != nil {
		reply.Data = err.Error()
		l.Data["json"] = &reply
		l.ServeJSON()
		return
	}
	queryInfo := queryData["guild_info"].(map[string]interface{})
	gameId := int(queryInfo["game_id"].(float64))
	logType := int(queryInfo["log_type"].(float64))
	serverId, _ := strconv.Atoi(queryInfo["server_id"].(string))
	guildName := queryInfo["guild_name"].(string)
	startTime, _ := strconv.ParseInt(strconv.FormatFloat(queryInfo["start_time"].(float64), 'f', -1, 64), 10, 64)
	endTime, _ := strconv.ParseInt(strconv.FormatFloat(queryInfo["end_time"].(float64), 'f', -1, 64), 10, 64)
	page := int(queryInfo["page"].(float64))
	pageSize := int(queryInfo["page_size"].(float64))
	recordName := queryInfo["record_name"].(string)

	query := GuildGameLogs{GameId: gameId, ServerId: serverId, LogType: logType, GuildName: guildName, StartTime: time.Unix(startTime/1000, 0).Format("2006-01-02 15:04:05"), EndTime: time.Unix(endTime/1000, 0).Format("2006-01-02 15:04:05"), Page: page, PageSize: pageSize}
	reply = QueryGuildGameLogs(&query, recordName)
	reply.Status = 200
	l.Data["json"] = &reply
	l.ServeJSON()

}
