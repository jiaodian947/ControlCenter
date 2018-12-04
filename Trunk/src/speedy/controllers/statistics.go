package controllers

import (
	"encoding/json"
	"speedy/server"
	"strconv"
	"time"
)

type Statistics struct {
	CheckRight
}

func (l *Statistics) GetGameStatistics() {
	var queryData map[string]interface{}
	var reply Reply
	reply.Status = 500
	if err := json.Unmarshal(l.Ctx.Input.RequestBody, &queryData); err != nil {
		reply.Data = err.Error()
		l.Data["json"] = &reply
		l.ServeJSON()
		return
	}
	data := queryData["data"].(map[string]interface{})
	gameId := int(data["game_id"].(float64))
	statisticsType := int(data["statistics_type"].(float64))
	serverId, _ := strconv.Atoi(data["server_id"].(string))
	startTime, _ := strconv.ParseInt(strconv.FormatFloat(data["start_time"].(float64), 'f', -1, 64), 10, 64)
	endTime, _ := strconv.ParseInt(strconv.FormatFloat(data["end_time"].(float64), 'f', -1, 64), 10, 64)
	recordName := data["record_name"].(string)

	query := QueryStatistics{GameId: gameId, ServerId: serverId, StatisticsType: statisticsType, StartTime: time.Unix(startTime/1000, 0).Format("2006-01-02 15:04:05"), EndTime: time.Unix(endTime/1000, 0).Format("2006-01-02 15:04:05")}
	reply = QueryGameStatistics(&query, recordName)
	reply.Status = 200
	l.Data["json"] = &reply
	l.ServeJSON()

}

func (l *Statistics) GetCapitalStatistics() {
	var queryData map[string]interface{}
	var reply Reply
	reply.Status = 500
	if err := json.Unmarshal(l.Ctx.Input.RequestBody, &queryData); err != nil {
		reply.Data = err.Error()
		l.Data["json"] = &reply
		l.ServeJSON()
		return
	}
	data := queryData["data"].(map[string]interface{})
	gameId := int(data["game_id"].(float64))
	statisticsType := int(data["statistics_type"].(float64))
	serverId, _ := strconv.Atoi(data["server_id"].(string))
	startTime, _ := strconv.ParseInt(strconv.FormatFloat(data["start_time"].(float64), 'f', -1, 64), 10, 64)
	endTime, _ := strconv.ParseInt(strconv.FormatFloat(data["end_time"].(float64), 'f', -1, 64), 10, 64)
	recordName := data["record_name"].(string)

	query := QueryStatistics{GameId: gameId, ServerId: serverId, StatisticsType: statisticsType, StartTime: time.Unix(startTime/1000, 0).Format("2006-01-02 15:04:05"), EndTime: time.Unix(endTime/1000, 0).Format("2006-01-02 15:04:05")}
	reply = QueryCapitalStatistics(&query, recordName)
	reply.Status = 200
	l.Data["json"] = &reply
	l.ServeJSON()

}
func (l *Statistics) GetCapitalStatisticsDetails() {
	var queryData map[string]interface{}
	var reply Reply
	reply.Status = 500
	if err := json.Unmarshal(l.Ctx.Input.RequestBody, &queryData); err != nil {
		reply.Data = err.Error()
		l.Data["json"] = &reply
		l.ServeJSON()
		return
	}
	data := queryData["data"].(map[string]interface{})
	gameId := int(data["game_id"].(float64))
	statisticsType := int(data["statistics_type"].(float64))
	capitalType, _ := strconv.Atoi(data["capital_type"].(string))
	serverId, _ := strconv.Atoi(data["server_id"].(string))
	query := CapitalStatisticsDetails{GameId: gameId, ServerId: serverId, StatisticsType: statisticsType, CapitalType: capitalType}
	reply = QueryCapitalStatisticsDetails(&query)
	reply.Status = 200
	l.Data["json"] = &reply
	l.ServeJSON()

}
func (l *Statistics) GetQuestionNaireStatistics() {
	var reply Reply
	reply.Status = 500
	serverId, _ := l.GetInt("server_id", 0)
	gameId, _ := l.GetInt("game_id", 0)
	page, _ := l.GetInt64("page", 0)
	pageSize, _ := l.GetInt64("pageSize", 10)
	gameServer := server.FindServerByServerId(gameId, serverId)
	if gameServer == nil || gameServer.DB == nil {
		reply.Data = "server not found"
		l.Data["json"] = &reply
		l.ServeJSON()
		return
	}
	// rows, err := gameServer.DB.Query("SELECT tb_config.s_id,tb_config.s_name,tb_config.n_type,tb_config.s_prop,tb_config.n_status,tb_config.n_version,tb_config.n_operid FROM tb_config limit ?,?", page*10, pageSize)
	rows, err := gameServer.DB.Query("SELECT * FROM tb_config where n_type = 3 limit ?,? ", page*pageSize, pageSize)

	if err != nil {
		reply.Data = err.Error()
		l.Data["json"] = &reply
		l.ServeJSON()
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
		"count": len(result),
		"data":  result,
	}
	l.Data["json"] = &reply
	l.ServeJSON()
	// fmt.Println(result)
	return
}
func (l *Statistics) GetQuestionNaireStatisticsDetails() {
	var queryData map[string]interface{}
	var reply Reply
	reply.Status = 500
	if err := json.Unmarshal(l.Ctx.Input.RequestBody, &queryData); err != nil {
		reply.Data = err.Error()
		l.Data["json"] = &reply
		l.ServeJSON()
		return
	}
	data := queryData["data"].(map[string]interface{})
	gameId := int(data["game_id"].(float64))
	questId, _ := strconv.Atoi(data["quest_id"].(string))
	answerId := data["answer_id"].([]interface{})
	serverId, _ := strconv.Atoi(data["server_id"].(string))
	query := QuestionNaireStatisticsDetails{GameId: gameId, ServerId: serverId, QuestId: questId, AnswerId: answerId}
	reply = QueryQuestionNaireStatisticsDetails(&query)
	reply.Status = 200
	l.Data["json"] = &reply
	l.ServeJSON()

}
func (l *Statistics) GetQuestionNaireNotMustStatisticsDetails() {
	var queryData map[string]interface{}
	var reply Reply
	reply.Status = 500
	if err := json.Unmarshal(l.Ctx.Input.RequestBody, &queryData); err != nil {
		reply.Data = err.Error()
		l.Data["json"] = &reply
		l.ServeJSON()
		return
	}
	data := queryData["data"].(map[string]interface{})
	gameId := int(data["game_id"].(float64))
	questId, _ := strconv.Atoi(data["quest_id"].(string))
	serverId, _ := strconv.Atoi(data["server_id"].(string))
	query := QuestionNaireNotMustStatisticsDetails{GameId: gameId, ServerId: serverId, QuestId: questId}
	reply = QueryQuestionNaireNotMustStatisticsDetails(&query)
	reply.Status = 200
	l.Data["json"] = &reply
	l.ServeJSON()

}
