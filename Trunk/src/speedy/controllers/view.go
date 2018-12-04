package controllers

import (
	"encoding/json"
	"speedy/models"
	"strconv"
	"time"
)

type View struct {
	CheckRight
}

func (v *View) GetViewList() {
	var err error
	var reply Reply
	reply.Status = 500
	views, err := models.GetAllViews(nil, nil, nil, nil, 0, 0)
	if err != nil {
		reply.Data = err.Error()
		v.Data["json"] = &reply
		v.ServeJSON()
		return
	}
	reply.Status = 200
	reply.Data = map[string]interface{}{
		"views": views,
	}
	v.Data["json"] = &reply
	v.ServeJSON()
}
func (v *View) GetViewByViewId() {
	var reply Reply
	reply.Status = 500
	viewId := v.GetString("view_id")
	view, err := models.GetAllViews(map[string]string{"view_id": viewId}, nil, nil, nil, 0, 0)
	if err != nil {
		reply.Data = err.Error()
		v.Data["json"] = &reply
		v.ServeJSON()
		return
	}
	reply.Status = 200
	reply.Data = map[string]interface{}{
		"viewDetails": view,
	}
	v.Data["json"] = &reply
	v.ServeJSON()
}
func (v *View) AddView() {
	var viewData map[string]interface{}
	var reply Reply
	reply.Status = 500

	if err := json.Unmarshal(v.Ctx.Input.RequestBody, &viewData); err != nil {
		reply.Data = err.Error()
		v.Data["json"] = &reply
		v.ServeJSON()
		return
	}
	viewDataItem, _ := viewData["viewData"].(map[string]interface{})
	data, _ := viewDataItem["data"].(map[string]interface{})
	createTime, _ := strconv.ParseInt(strconv.FormatFloat(data["create_time"].(float64), 'f', -1, 64), 10, 64)
	title := data["title"].(string)
	path := data["path"].(string)
	category := int(data["category"].(float64))
	var categoryDesc string
	switch category {
	case 1:
		categoryDesc = "邮件管理"
	case 2:
		categoryDesc = "账户管理"
	case 3:
		categoryDesc = "活动管理"
	case 4:
		categoryDesc = "后台设置"
	case 5:
		categoryDesc = "玩家管理"
	case 6:
		categoryDesc = "公共数据"
	}
	hidden := 0
	if data["hidden"].(bool) {
		hidden = 1
	}
	viewItem := models.Views{CreateTime: time.Unix(createTime/1000, 0), Title: title, Path: path, Category: category, CategoryDesc: categoryDesc, Hidden: hidden}
	_, err := models.AddViews(&viewItem)
	if err != nil {
		reply.Data = err.Error()
		v.Data["json"] = &reply
		v.ServeJSON()
		return
	}
	reply.Status = 200
	v.Data["json"] = &reply
	v.ServeJSON()
	return
}

// func SmartPrint(i interface{}) {
// 	var kv = make(map[string]interface{})
// 	vValue := reflect.ValueOf(i)
// 	vType := reflect.TypeOf(i)
// 	for i := 0; i < vValue.NumField(); i++ {
// 		kv[vType.Field(i).Name] = vValue.Field(i)
// 	}
// 	fmt.Println("获取到数据:")
// 	for k, v := range kv {
// 		fmt.Print(k)
// 		fmt.Print(":")
// 		fmt.Print(v)
// 		fmt.Println()
// 	}
// }
