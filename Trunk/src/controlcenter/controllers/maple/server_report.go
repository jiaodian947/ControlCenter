package maple

import (
	"controlcenter/modules/models"
	"encoding/json"
	"log"

	"github.com/astaxie/beego"
)

type PlayerCount struct {
	GameId     int
	DistrictId int
	ServerId   int
	Players    int
	Capacity   int
}

func ReportServerStatus(handler *Handler) {
	var pc PlayerCount
	err := json.Unmarshal(handler.data, &pc)
	if err != nil {
		log.Println(err)
		return
	}

	var server models.ServerInfo
	server.GameId = pc.GameId
	server.Id = pc.ServerId
	server.DistrictId = pc.DistrictId
	if err := server.Read("game_id", "id", "district_id"); err != nil {
		beego.Warning("server not found", server)
		return
	}

	server.PlayerCount = pc.Players
	server.PlayerMaxCount = pc.Capacity
	server.Update("player_count", "player_max_count")
	beego.Info("update server info")
}
