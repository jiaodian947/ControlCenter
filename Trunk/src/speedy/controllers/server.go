package controllers

import (
	"encoding/json"
	"speedy/models"
	"strconv"
)

type Server struct {
	CheckRight
}

func (s *Server) GetServerList() {
	var err error
	var reply Reply
	reply.Status = 500
	servers, err := models.GetAllServer(nil, nil, nil, nil, 0, 0)
	if err != nil {
		reply.Data = err.Error()
		s.Data["json"] = &reply
		s.ServeJSON()
		return
	}
	reply.Status = 200
	reply.Data = map[string]interface{}{
		"servers": servers,
	}
	s.Data["json"] = &reply
	s.ServeJSON()
}
func (s *Server) GetDomainInfo() {
	var queryData map[string]interface{}
	var reply Reply
	reply.Status = 500
	if err := json.Unmarshal(s.Ctx.Input.RequestBody, &queryData); err != nil {
		reply.Data = err.Error()
		s.Data["json"] = &reply
		s.ServeJSON()
		return
	}
	queryInfo := queryData["server_info"].(map[string]interface{})
	gameId := int(queryInfo["game_id"].(float64))
	serverId, _ := strconv.Atoi(queryInfo["server_id"].(string))
	recordName := queryInfo["record_name"].(string)
	query := DoMainsQueryRecord{GameId: gameId, ServerId: serverId, RecordName: recordName}
	reply = QueryDoMain(&query)

	reply.Status = 200
	s.Data["json"] = &reply
	s.ServeJSON()
}
func (s *Server) GetDomainRecordName() {
	var queryData map[string]interface{}
	var reply Reply
	reply.Status = 500
	if err := json.Unmarshal(s.Ctx.Input.RequestBody, &queryData); err != nil {
		reply.Data = err.Error()
		s.Data["json"] = &reply
		s.ServeJSON()
		return
	}
	queryInfo := queryData["server_info"].(map[string]interface{})
	// type Query struct {
	// GameId   int    `json:"gameid"`
	// ServerId int    `json:"serverid"`
	// RoleName string `json:"rolename"`
	// }
	gameId := int(queryInfo["game_id"].(float64))
	serverId, _ := strconv.Atoi(queryInfo["server_id"].(string))
	query := Query{GameId: gameId, ServerId: serverId}
	reply = QueryDoMainRecordName(&query)
	reply.Status = 200
	s.Data["json"] = &reply
	s.ServeJSON()
}
