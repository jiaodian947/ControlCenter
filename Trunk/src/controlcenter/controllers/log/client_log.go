package log

import (
	"controlcenter/controllers"
	"controlcenter/modules/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"sort"

	"github.com/astaxie/beego"
)

type LogResult struct {
	Ok bool
}

type ClientLog struct {
	controllers.BaseRouter
}

func (c *ClientLog) GetUpload() {
	c.TplName = "log_upload.html"
}

func (c *ClientLog) Upload() {
	params := c.GetString("params")

	for k, _ := range c.Ctx.Request.MultipartForm.File {
		f, h, err := c.GetFile(k)
		if err != nil {
			c.Data["json"] = &LogResult{false}
			c.ServeJSON()
			return
		}
		// 获取当前年月
		datename := time.Now().Format("2006_01_02_15_04_05")
		// 设置保存目录
		dirPath := "./static/upload/"
		// 设置保存文件名
		FileName := h.Filename

		path := fmt.Sprintf("%s/%s%s%s", dirPath, datename, utils.GetRandomString(4), FileName)
		f.Close()
		err = c.SaveToFile(k, path)
		if err != nil {
			c.Data["json"] = &LogResult{false}
			c.ServeJSON()
			return
		}

		ioutil.WriteFile(path+".json", []byte(params), 0666)
	}

	c.Data["json"] = &LogResult{true}
	c.ServeJSON()

}

type FileInfo struct {
	LogInfo    *LogInfo
	FileUrl    string
	FileName   string
	CreateTime string
	Size       int64
}

type LogInfo struct {
	Version      string
	ProjectName  string
	RoleName     string
	Guid         string
	Device       string
	OsVer        string
	Ptime        string
	FMemory      string
	GameId       string
	Type         string
	LifeSpan     string
	SceneId      string
	UserAccount  string
	ChannelId    string
	SdkVersion   string
	TotalMemory  string
	IsRoot       string
	CrashPackage string
}

func (c *ClientLog) GetLogs() {
	c.TplName = "log_list.html"
	dir_list, err := ioutil.ReadDir("./static/upload")
	if err != nil {
		return
	}

	curpage, err := c.GetInt("p")
	if err != nil {
		curpage = 1
	}

	curpage--
	pagelimit := 20
	fis := make([]FileInfo, 0, pagelimit)

	log_files := make([]os.FileInfo, 0, len(dir_list))
	for _, v := range dir_list {
		if !strings.HasSuffix(v.Name(), ".json") {
			log_files = append(log_files, v)
		}
	}

	sort.Slice(log_files, func(i, j int) bool { return log_files[i].ModTime().After(log_files[j].ModTime()) })

	total := len(log_files)
	c.SetPaginator(pagelimit, int64(total))
	startpos := curpage * pagelimit
	stoppos := startpos + pagelimit
	for k, v := range log_files {
		if k < startpos {
			continue
		}
		if k >= stoppos {
			break
		}

		fi := FileInfo{}

		fi.FileUrl = fmt.Sprintf("/static/upload/%s", v.Name())
		info, err := ioutil.ReadFile(fmt.Sprintf("./static/upload/%s.json", v.Name()))
		if err == nil {
			logInfo := &LogInfo{}
			if err := json.Unmarshal(info, logInfo); err == nil {
				fi.LogInfo = logInfo
			} else {
				beego.Error(err)
			}
		}

		fi.FileName = v.Name()
		fi.CreateTime = v.ModTime().Format("2006/01/02 15:04:05")
		fi.Size = v.Size()
		fis = append(fis, fi)
	}
	c.Data["files"] = fis

}
