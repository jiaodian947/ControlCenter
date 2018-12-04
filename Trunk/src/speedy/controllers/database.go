package controllers

import (
	"fmt"
	"speedy/models/gameobj"
	"speedy/server"
)

type Record map[string]interface{}

// type RecordName string

// const (
// 	CAPITALLOG     RecordName = "role_capital_log"
// 	RANKLISTLOG               = "role_ranklist_log"
// 	ACHIEVEMENTLOG            = "role_achievement_log"
// 	TITLELOG                  = "role_title_log"
// )

type DatabaseController struct {
	BaseRouter
}

type PlayerRoles struct {
	RoleId   int64
	RoleName string
	Passport string
	SaveData []byte
}
type DoMainsQueryRecord struct {
	GameId     int    `json:"gameid"`
	ServerId   int    `json:"serverid"`
	RecordName string `json:"recordname"`
}
type DoMains struct {
	RecordName string
	SaveData   []byte
}

type Query struct {
	GameId   int    `json:"gameid"`
	ServerId int    `json:"serverid"`
	RoleName string `json:"rolename"`
}
type QueryLogs struct {
	RecordName string
	GameId     int
	ServerId   int
	RoleName   string
	StartTime  string
	EndTime    string
	Page       int
	PageSize   int
}
type GuildFunctionLogs struct {
	RecordName string
	GameId     int
	ServerId   int
	GuildName  string
	StartTime  string
	EndTime    string
	Page       int
	PageSize   int
}
type GuildGameLogs struct {
	RecordName string
	GameId     int
	ServerId   int
	LogType    int
	GuildName  string
	StartTime  string
	EndTime    string
	Page       int
	PageSize   int
}

type QueryStatistics struct {
	RecordName     string
	GameId         int
	ServerId       int
	StatisticsType int
	StartTime      string
	EndTime        string
}
type CapitalStatisticsDetails struct {
	GameId         int
	ServerId       int
	StatisticsType int
	CapitalType    int
}

type QuestionNaireStatisticsDetails struct {
	GameId   int
	ServerId int
	QuestId  int
	AnswerId []interface{}
}
type QuestionNaireNotMustStatisticsDetails struct {
	GameId   int
	ServerId int
	QuestId  int
}

func QueryPlayerInfo(query *Query) (res Reply) {
	res.Status = 500

	defer func() {
		if p := recover(); p != nil {
			switch inst := p.(type) {
			case error:
				res.Data = inst.Error()
			case string:
				res.Data = inst
			}
			res.Status = 500

		}
	}()

	gs := server.FindServerByServerId(query.GameId, query.ServerId)
	if gs == nil || gs.DB == nil {
		res.Data = "server not found"
		return
	}

	row, err := gs.DB.Query("select r.n_roleid, r.s_rolename, r.s_passport, b.lb_save_data from player_roles as r, player_binary as b where r.n_roleid=b.n_roleid and r.s_rolename = ?", query.RoleName)

	if err != nil {
		res.Data = err.Error()
		return
	}

	defer row.Close()

	r := &PlayerRoles{}
	if row.Next() {
		err := row.Scan(&r.RoleId, &r.RoleName, &r.Passport, &r.SaveData)
		if err != nil {
			res.Data = err.Error()
			return
		}
	}

	if len(r.SaveData) == 0 {
		res.Data = "role not found"
		return
	}
	obj := gameobj.NewGameObjectFromBinary(r.SaveData)
	res.Status = 200
	obj.Account = r.Passport
	res.Data = obj
	return
}
func QueryDoMain(query *DoMainsQueryRecord) (res Reply) {
	res.Status = 500

	defer func() {
		if p := recover(); p != nil {
			switch inst := p.(type) {
			case error:
				res.Data = inst.Error()
			case string:
				res.Data = inst
			}
			res.Status = 500

		}
	}()

	gs := server.FindServerByServerId(query.GameId, query.ServerId)
	if gs == nil || gs.DB == nil {
		res.Data = "server not found"
		return
	}
	row, err := gs.DB.Query("SELECT lb_save_data from tb_domains WHERE s_name = ?", query.RecordName)
	// row, err := gs.DB.Query(str)

	if err != nil {
		res.Data = err.Error()
		return
	}

	defer row.Close()

	d := &DoMains{}
	if row.Next() {
		err := row.Scan(&d.SaveData)
		if err != nil {
			res.Data = err.Error()
			return
		}
	}

	if len(d.SaveData) == 0 {
		res.Data = "role not found"
		return
	}
	obj := gameobj.NewGameDataFromBinary(d.SaveData)
	res.Status = 200
	res.Data = obj
	return
}
func QueryDoMainRecordName(query *Query) (res Reply) {
	res.Status = 500

	defer func() {
		if p := recover(); p != nil {
			switch inst := p.(type) {
			case error:
				res.Data = inst.Error()
			case string:
				res.Data = inst
			}
			res.Status = 500

		}
	}()

	gs := server.FindServerByServerId(query.GameId, query.ServerId)
	if gs == nil || gs.DB == nil {
		res.Data = "server not found"
		return
	}
	str := "SELECT s_name from tb_domains"
	// row, err := gs.DB.Query("SELECT s_name, lb_save_data from tb_domains WHERE s_name LIKE '%_?' ESCAPE '_' ", query.ServerId)
	row, err := gs.DB.Query(str)

	if err != nil {
		res.Data = err.Error()
		return
	}

	defer row.Close()

	d := &DoMains{}
	var rowsArr []string
	for row.Next() {
		err := row.Scan(&d.RecordName)
		if err != nil {
			res.Data = err.Error()
			return
		}
		rowsArr = append(rowsArr, d.RecordName)
	}

	if len(rowsArr) == 0 {
		res.Data = "role not found"
		return
	}
	// obj := gameobj.NewGameDataFromBinary(d.SaveData)
	res.Status = 200
	res.Data = rowsArr
	return
}

func QueryGameLogs(query *QueryLogs, recordName string) (res Reply) {
	res.Status = 500
	defer func() {
		if p := recover(); p != nil {
			switch inst := p.(type) {
			case error:
				res.Data = inst.Error()
			case string:
				res.Data = inst
			}
			res.Status = 500

		}
	}()

	gs := server.FindServerByServerId(query.GameId, query.ServerId)

	if gs == nil || gs.LogDB == nil {
		res.Data = "server not found"
		return
	}
	if gs == nil || gs.DB == nil {
		res.Data = "server not found"
		return
	}
	// str := "SELECT role_name,capital_type,value,item_log_type,comment_1,comment_2,comment_3 from role_capital_log from role_capital_log where role_name = ?"
	// row, err := gs.DB.Query("SELECT s_name, lb_save_data from tb_domains WHERE s_name LIKE '%_?' ESCAPE '_' ", query.ServerId)
	var NRoleid int64
	dbrows, dberr := gs.DB.Query("select n_roleid from player_roles where s_rolename = ?", query.RoleName)
	if dberr != nil {
		res.Data = dberr.Error()
		return
	}
	defer dbrows.Close()
	for dbrows.Next() {
		err := dbrows.Scan(&NRoleid)
		if err != nil {
			fmt.Println(err)
		}
	}
	sqlStr := fmt.Sprintf("SELECT * from %s where role_id = %d and  DATE_FORMAT(change_time,'%%Y-%%m-%%d %%H:%%i:%%s') > '%s' and DATE_FORMAT(change_time,'%%Y-%%m-%%d %%H:%%i:%%s') < '%s' limit %d,%d", recordName, NRoleid, query.StartTime, query.EndTime, query.Page*query.PageSize, query.PageSize)
	fmt.Print(sqlStr)
	row, err := gs.LogDB.Query(sqlStr)

	if err != nil {
		res.Data = err.Error()
		return
	}

	defer row.Close()

	var rowsArr []Record
	columns, err := row.Columns()
	if err != nil {
		res.Data = err.Error()
		return
	}

	for row.Next() {
		r := make([]interface{}, len(columns))
		container := make([]interface{}, len(columns))
		for i := range r {
			r[i] = &container[i]
		}
		err := row.Scan(r...)
		if err != nil {
			res.Data = err.Error()
			return
		}
		var record Record = make(map[string]interface{}, len(columns))
		for i, column := range columns {
			record[column] = container[i]
		}
		for k, v := range record {
			switch v.(type) {
			case []uint8:
				arr := v.([]uint8)
				record[k] = uiToS(arr)
			case nil:
				record[k] = ""
			}
		}

		rowsArr = append(rowsArr, record)
	}

	if len(rowsArr) == 0 {
		res.Data = "role not found"
		return
	}
	// obj := gameobj.NewGameDataFromBinary(d.SaveData)
	res.Status = 200
	res.Data = rowsArr
	return
}
func QueryGuildFunctionLogs(query *GuildFunctionLogs, recordName string) (res Reply) {
	res.Status = 500
	defer func() {
		if p := recover(); p != nil {
			switch inst := p.(type) {
			case error:
				res.Data = inst.Error()
			case string:
				res.Data = inst
			}
			res.Status = 500

		}
	}()

	gs := server.FindServerByServerId(query.GameId, query.ServerId)

	if gs == nil || gs.LogDB == nil {
		res.Data = "server not found"
		return
	}

	sqlStr := fmt.Sprintf("SELECT * from %s where guild_name = '%s' and  DATE_FORMAT(change_time,'%%Y-%%m-%%d %%H:%%i:%%s') > '%s' and DATE_FORMAT(change_time,'%%Y-%%m-%%d %%H:%%i:%%s') < '%s' limit %d,%d", recordName, query.GuildName, query.StartTime, query.EndTime, query.Page*query.PageSize, query.PageSize)
	fmt.Print(sqlStr)
	row, err := gs.LogDB.Query(sqlStr)

	if err != nil {
		res.Data = err.Error()
		return
	}

	defer row.Close()

	var rowsArr []Record
	columns, err := row.Columns()
	if err != nil {
		res.Data = err.Error()
		return
	}

	for row.Next() {
		r := make([]interface{}, len(columns))
		container := make([]interface{}, len(columns))
		for i := range r {
			r[i] = &container[i]
		}
		err := row.Scan(r...)
		if err != nil {
			res.Data = err.Error()
			return
		}
		var record Record = make(map[string]interface{}, len(columns))
		for i, column := range columns {
			record[column] = container[i]
		}
		for k, v := range record {
			switch v.(type) {
			case []uint8:
				arr := v.([]uint8)
				record[k] = uiToS(arr)
			case nil:
				record[k] = ""
			}
		}

		rowsArr = append(rowsArr, record)
	}

	if len(rowsArr) == 0 {
		res.Data = "role not found"
		return
	}
	// obj := gameobj.NewGameDataFromBinary(d.SaveData)
	res.Status = 200
	res.Data = rowsArr
	return
}
func QueryGuildGameLogs(query *GuildGameLogs, recordName string) (res Reply) {
	res.Status = 500
	defer func() {
		if p := recover(); p != nil {
			switch inst := p.(type) {
			case error:
				res.Data = inst.Error()
			case string:
				res.Data = inst
			}
			res.Status = 500

		}
	}()

	gs := server.FindServerByServerId(query.GameId, query.ServerId)

	if gs == nil || gs.LogDB == nil {
		res.Data = "server not found"
		return
	}

	sqlStr := fmt.Sprintf("SELECT * from %s where guild_name = '%s' and log_type = %d and  DATE_FORMAT(change_time,'%%Y-%%m-%%d %%H:%%i:%%s') > '%s' and DATE_FORMAT(change_time,'%%Y-%%m-%%d %%H:%%i:%%s') < '%s' limit %d,%d", recordName, query.GuildName, query.LogType, query.StartTime, query.EndTime, query.Page*query.PageSize, query.PageSize)
	fmt.Print(sqlStr)
	row, err := gs.LogDB.Query(sqlStr)

	if err != nil {
		res.Data = err.Error()
		return
	}

	defer row.Close()

	var rowsArr []Record
	columns, err := row.Columns()
	if err != nil {
		res.Data = err.Error()
		return
	}

	for row.Next() {
		r := make([]interface{}, len(columns))
		container := make([]interface{}, len(columns))
		for i := range r {
			r[i] = &container[i]
		}
		err := row.Scan(r...)
		if err != nil {
			res.Data = err.Error()
			return
		}
		var record Record = make(map[string]interface{}, len(columns))
		for i, column := range columns {
			record[column] = container[i]
		}
		for k, v := range record {
			switch v.(type) {
			case []uint8:
				arr := v.([]uint8)
				record[k] = uiToS(arr)
			case nil:
				record[k] = ""
			}
		}

		rowsArr = append(rowsArr, record)
	}

	if len(rowsArr) == 0 {
		res.Data = "role not found"
		return
	}
	// obj := gameobj.NewGameDataFromBinary(d.SaveData)
	res.Status = 200
	res.Data = rowsArr
	return
}

//Statistics

// CapitalStatistics
func QueryGameStatistics(query *QueryStatistics, recordName string) (res Reply) {
	res.Status = 500
	defer func() {
		if p := recover(); p != nil {
			switch inst := p.(type) {
			case error:
				res.Data = inst.Error()
			case string:
				res.Data = inst
			}
			res.Status = 500

		}
	}()

	gs := server.FindServerByServerId(query.GameId, query.ServerId)

	if gs == nil || gs.LogDB == nil {
		res.Data = "server not found"
		return
	}

	var sqlStr string
	if query.StatisticsType == 0 {
		sqlStr = fmt.Sprintf("SELECT * from %s where log_type < 40000 and  DATE_FORMAT(change_time,'%%Y-%%m-%%d %%H:%%i:%%s') > '%s' and DATE_FORMAT(change_time,'%%Y-%%m-%%d %%H:%%i:%%s') < '%s'", recordName, query.StartTime, query.EndTime)

	}
	if query.StatisticsType == 1 {
		sqlStr = fmt.Sprintf("SELECT * from %s where log_type > 40000 and  DATE_FORMAT(change_time,'%%Y-%%m-%%d %%H:%%i:%%s') > '%s' and DATE_FORMAT(change_time,'%%Y-%%m-%%d %%H:%%i:%%s') < '%s'", recordName, query.StartTime, query.EndTime)

	}
	row, err := gs.LogDB.Query(sqlStr)

	if err != nil {
		res.Data = err.Error()
		return
	}

	defer row.Close()

	var rowsArr []Record
	columns, err := row.Columns()
	if err != nil {
		res.Data = err.Error()
		return
	}

	for row.Next() {
		r := make([]interface{}, len(columns))
		container := make([]interface{}, len(columns))
		for i := range r {
			r[i] = &container[i]
		}
		err := row.Scan(r...)
		if err != nil {
			res.Data = err.Error()
			return
		}
		var record Record = make(map[string]interface{}, len(columns))
		for i, column := range columns {
			record[column] = container[i]
		}
		for k, v := range record {
			switch v.(type) {
			case []uint8:
				arr := v.([]uint8)
				record[k] = uiToS(arr)
			case nil:
				record[k] = ""
			}
		}

		rowsArr = append(rowsArr, record)
	}

	if len(rowsArr) == 0 {
		res.Data = "role not found"
		return
	}
	// obj := gameobj.NewGameDataFromBinary(d.SaveData)
	res.Status = 200
	res.Data = rowsArr
	return
}
func QueryCapitalStatistics(query *QueryStatistics, recordName string) (res Reply) {
	res.Status = 500
	defer func() {
		if p := recover(); p != nil {
			switch inst := p.(type) {
			case error:
				res.Data = inst.Error()
			case string:
				res.Data = inst
			}
			res.Status = 500

		}
	}()

	gs := server.FindServerByServerId(query.GameId, query.ServerId)

	if gs == nil || gs.LogDB == nil {
		res.Data = "server not found"
		return
	}

	var sqlStr string
	if query.StatisticsType == 0 {
		sqlStr = fmt.Sprintf("select capital_type,SUM(value) as sum,SUM(value)/(COUNT(DISTINCT role_id)) as avg from %s where log_type < 40000 and  DATE_FORMAT(change_time,'%%Y-%%m-%%d %%H:%%i:%%s') > '%s' and DATE_FORMAT(change_time,'%%Y-%%m-%%d %%H:%%i:%%s') < '%s' group by capital_type", recordName, query.StartTime, query.EndTime)

	}
	if query.StatisticsType == 1 {
		sqlStr = fmt.Sprintf("select capital_type,SUM(value) as sum,SUM(value)/(COUNT(DISTINCT role_id)) as avg from %s where log_type > 40000 and  DATE_FORMAT(change_time,'%%Y-%%m-%%d %%H:%%i:%%s') > '%s' and DATE_FORMAT(change_time,'%%Y-%%m-%%d %%H:%%i:%%s') < '%s'  group by capital_type", recordName, query.StartTime, query.EndTime)

	}
	row, err := gs.LogDB.Query(sqlStr)

	if err != nil {
		res.Data = err.Error()
		return
	}

	defer row.Close()

	var rowsArr []Record
	columns, err := row.Columns()
	if err != nil {
		res.Data = err.Error()
		return
	}

	for row.Next() {
		r := make([]interface{}, len(columns))
		container := make([]interface{}, len(columns))
		for i := range r {
			r[i] = &container[i]
		}
		err := row.Scan(r...)
		if err != nil {
			res.Data = err.Error()
			return
		}
		var record Record = make(map[string]interface{}, len(columns))
		for i, column := range columns {
			record[column] = container[i]
		}
		for k, v := range record {
			switch v.(type) {
			case []uint8:
				arr := v.([]uint8)
				record[k] = uiToS(arr)
			case nil:
				record[k] = ""
			}
		}

		rowsArr = append(rowsArr, record)
	}

	if len(rowsArr) == 0 {
		res.Data = "role not found"
		return
	}
	// obj := gameobj.NewGameDataFromBinary(d.SaveData)
	res.Status = 200
	res.Data = rowsArr
	return
}
func QueryCapitalStatisticsDetails(query *CapitalStatisticsDetails) (res Reply) {
	res.Status = 500
	defer func() {
		if p := recover(); p != nil {
			switch inst := p.(type) {
			case error:
				res.Data = inst.Error()
			case string:
				res.Data = inst
			}
			res.Status = 500

		}
	}()

	gs := server.FindServerByServerId(query.GameId, query.ServerId)

	if gs == nil || gs.LogDB == nil {
		res.Data = "server not found"
		return
	}

	var sqlStr string
	if query.StatisticsType == 0 {
		sqlStr = fmt.Sprintf("select DISTINCT log_type,SUM(value) as sum from role_capital_log where log_type < 40000 and capital_type = %d group by log_type ", query.CapitalType)

	}
	if query.StatisticsType == 1 {
		sqlStr = fmt.Sprintf("select DISTINCT log_type,SUM(value) as sum from role_capital_log where log_type > 40000 and capital_type = %d group by log_type ", query.CapitalType)

	}
	row, err := gs.LogDB.Query(sqlStr)

	if err != nil {
		res.Data = err.Error()
		return
	}

	defer row.Close()

	var rowsArr []Record
	columns, err := row.Columns()
	if err != nil {
		res.Data = err.Error()
		return
	}

	for row.Next() {
		r := make([]interface{}, len(columns))
		container := make([]interface{}, len(columns))
		for i := range r {
			r[i] = &container[i]
		}
		err := row.Scan(r...)
		if err != nil {
			res.Data = err.Error()
			return
		}
		var record Record = make(map[string]interface{}, len(columns))
		for i, column := range columns {
			record[column] = container[i]
		}
		for k, v := range record {
			switch v.(type) {
			case []uint8:
				arr := v.([]uint8)
				record[k] = uiToS(arr)
			case nil:
				record[k] = ""
			}
		}

		rowsArr = append(rowsArr, record)
	}

	if len(rowsArr) == 0 {
		res.Data = "role not found"
		return
	}
	// obj := gameobj.NewGameDataFromBinary(d.SaveData)
	res.Status = 200
	res.Data = rowsArr
	return
}

func QueryQuestionNaireStatisticsDetails(query *QuestionNaireStatisticsDetails) (res Reply) {
	res.Status = 500
	defer func() {
		if p := recover(); p != nil {
			switch inst := p.(type) {
			case error:
				res.Data = inst.Error()
			case string:
				res.Data = inst
			}
			res.Status = 500

		}
	}()

	gs := server.FindServerByServerId(query.GameId, query.ServerId)

	if gs == nil || gs.LogDB == nil {
		res.Data = "server not found"
		return
	}
	var rowsArr []Record

	for i := 0; i < len(query.AnswerId); i++ {
		var sqlStr string
		sqlStr = fmt.Sprintf("SELECT %s as answer_id,COUNT(*) as count FROM questionnaire_log WHERE quest_id =%d  AND version = (SELECT MAX(version)  FROM questionnaire_log)  AND answer like '%%%s%%' ", query.AnswerId[i], query.QuestId, query.AnswerId[i])
		row, err := gs.LogDB.Query(sqlStr)

		if err != nil {
			res.Data = err.Error()
			return
		}

		defer row.Close()

		columns, err := row.Columns()
		if err != nil {
			res.Data = err.Error()
			return
		}

		for row.Next() {
			r := make([]interface{}, len(columns))
			container := make([]interface{}, len(columns))
			for i := range r {
				r[i] = &container[i]
			}
			err := row.Scan(r...)
			if err != nil {
				res.Data = err.Error()
				return
			}
			var record Record = make(map[string]interface{}, len(columns))
			for i, column := range columns {
				record[column] = container[i]
			}
			for k, v := range record {
				switch v.(type) {
				case []uint8:
					arr := v.([]uint8)
					record[k] = uiToS(arr)
				case nil:
					record[k] = ""
				}
			}

			rowsArr = append(rowsArr, record)

		}
	}

	if len(rowsArr) == 0 {
		res.Data = "role not found"
		return
	}
	// obj := gameobj.NewGameDataFromBinary(d.SaveData)
	res.Status = 200
	res.Data = rowsArr
	return
}
func QueryQuestionNaireNotMustStatisticsDetails(query *QuestionNaireNotMustStatisticsDetails) (res Reply) {
	res.Status = 500
	defer func() {
		if p := recover(); p != nil {
			switch inst := p.(type) {
			case error:
				res.Data = inst.Error()
			case string:
				res.Data = inst
			}
			res.Status = 500

		}
	}()

	gs := server.FindServerByServerId(query.GameId, query.ServerId)

	if gs == nil || gs.LogDB == nil {
		res.Data = "server not found"
		return
	}
	var rowsArr []Record

	var sqlStr string
	sqlStr = fmt.Sprintf("SELECT *,COUNT(*) as count FROM questionnaire_log WHERE quest_id =%d  AND version = (SELECT MAX(version)  FROM questionnaire_log)   ", query.QuestId)
	row, err := gs.LogDB.Query(sqlStr)

	if err != nil {
		res.Data = err.Error()
		return
	}

	defer row.Close()

	columns, err := row.Columns()
	if err != nil {
		res.Data = err.Error()
		return
	}

	for row.Next() {
		r := make([]interface{}, len(columns))
		container := make([]interface{}, len(columns))
		for i := range r {
			r[i] = &container[i]
		}
		err := row.Scan(r...)
		if err != nil {
			res.Data = err.Error()
			return
		}
		var record Record = make(map[string]interface{}, len(columns))
		for i, column := range columns {
			record[column] = container[i]
		}
		for k, v := range record {
			switch v.(type) {
			case []uint8:
				arr := v.([]uint8)
				record[k] = uiToS(arr)
			case nil:
				record[k] = ""
			}
		}

		rowsArr = append(rowsArr, record)

	}

	if len(rowsArr) == 0 {
		res.Data = "role not found"
		return
	}
	// obj := gameobj.NewGameDataFromBinary(d.SaveData)
	res.Status = 200
	res.Data = rowsArr
	return
}
func uiToS(bs []uint8) string {
	ba := []byte{}
	for _, b := range bs {
		ba = append(ba, byte(b))
	}
	return string(ba)
}
