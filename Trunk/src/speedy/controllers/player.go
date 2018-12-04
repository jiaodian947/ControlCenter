package controllers

import (
	"encoding/json"
	"strconv"
)

type Player struct {
	CheckRight
}

// type PlayerXml struct {
// 	XMLName    xml.Name     `xml:"object"`
// 	Properties []Properties `xml:"properties"`
// 	Records    []Records    `xml:"records"`
// }
// type Properties struct {
// 	XMLName  xml.Name   `xml:"properties"`
// 	Property []Property `xml:"property"`
// }
// type Property struct {
// 	XMLName xml.Name `xml:"property"`
// 	Name    string   `xml:"name,attr"`
// 	Desc    string   `xml:"desc,attr"`
// }
// type Records struct {
// 	XMLName xml.Name `xml:"records"`
// 	Record  []Record `xml:"record"`
// }
// type Record struct {
// 	XMLName xml.Name `xml:"record"`
// 	Name    string   `xml:"name,attr"`
// 	Desc    string   `xml:"desc,attr"`
// 	Column  []Column `xml:"column"`
// }
// type Column struct {
// 	XMLName xml.Name `xml:"column"`
// 	Desc    string   `xml:"desc,attr"`
// }

// func file_open(file_name string) ([]byte, error) {
// 	//定义变量
// 	var (
// 		open               *os.File
// 		file_data          []byte
// 		open_err, read_err error
// 	)
// 	//打开文件
// 	open, open_err = os.Open(file_name)
// 	if open_err != nil {
// 		return nil, open_err
// 	}
// 	//关闭资源
// 	defer open.Close()

// 	//读取所有文件内容
// 	file_data, read_err = ioutil.ReadAll(open)
// 	if read_err != nil {
// 		return nil, read_err
// 	}
// 	return file_data, nil
// }
func (p *Player) GetPlayerInfo() {
	var queryData map[string]interface{}
	var reply Reply
	reply.Status = 500
	if err := json.Unmarshal(p.Ctx.Input.RequestBody, &queryData); err != nil {
		reply.Data = err.Error()
		p.Data["json"] = &reply
		p.ServeJSON()
		return
	}
	queryInfo := queryData["player_info"].(map[string]interface{})
	// type Query struct {
	// GameId   int    `json:"gameid"`
	// ServerId int    `json:"serverid"`
	// RoleName string `json:"rolename"`
	// }
	gameId := int(queryInfo["game_id"].(float64))
	serverId, _ := strconv.Atoi(queryInfo["server_id"].(string))
	roleName := queryInfo["role_name"].(string)
	query := Query{GameId: gameId, ServerId: serverId, RoleName: roleName}
	reply = QueryPlayerInfo(&query)
	// file_data, err := file_open("player.xml")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	os.Exit(0)
	// }
	// tpl := PlayerXml{}
	// tpl_err := xml.Unmarshal(file_data, &tpl)
	// if tpl_err != nil {
	// 	fmt.Println("tpl err :", tpl_err)
	// 	os.Exit(0)
	// }
	// fmt.Println(tpl)
	// //替换玩家属性中的key字段

	// replyJson, _ := json.Marshal(reply)
	// var playerInfoData map[string]interface{}
	// if err := json.Unmarshal(replyJson, &playerInfoData); err != nil {
	// 	reply.Data = err.Error()
	// 	p.Data["json"] = &reply
	// 	p.ServeJSON()
	// 	return
	// }
	// playerAttrData := playerInfoData["Data"].(map[string]interface{})["Attrs"].(map[string]interface{})["Attr"].([]map[string]interface{})

	// for i := 0; i < len(playerAttrData); i++ {
	// 	for _, j := range tpl.Properties[0].Property {
	// 		if j.Name == playerAttrData[i]["Key"].(string) {
	// 			playerAttrData[i]["Val"] = j.Desc
	// 		}
	// 	}
	// }
	// fmt.Println(playerAttrData)

	reply.Status = 200
	p.Data["json"] = &reply
	p.ServeJSON()

}
