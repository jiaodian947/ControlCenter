package utils

import (
	"encoding/json"
	"io/ioutil"
)

var JsonConf = new(ConfJson)

//定义配置文件解析后的结构
type ConfJson struct {
	StartPort      int `json:"start_port"`
	MoveInterval   int `json:"move_interval"`
	AttackInterval int `json:"attack_interval"`
	ChatInterval   int `json:"chat_interval"`
}

type JsonStruct struct {
}

func NewJsonStruct() *JsonStruct {
	return &JsonStruct{}
}

func (jst *JsonStruct) Load(filename string, v interface{}) {
	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal(data, v)
	if err != nil {
		return
	}
}
