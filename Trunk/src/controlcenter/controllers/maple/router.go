package maple

import (
	"github.com/astaxie/beego"
)

func AddMapleRouter() {
	gc := new(GameController)
	beego.Router("/game", gc)
	beego.Router("/game/:id:int/del", gc, "get:DeleteGame")

	beego.Router("/game/add", &GameAddController{})

	dc := new(DistrictController)
	beego.Router("/game/:id:int", dc)
	beego.Router("/game/:id1:int/:id2:int/edit", dc, "get:EditDistrict;post:UpdateDistrict")
	beego.Router("/game/:id1:int/:id2:int/del", dc, "get:DeleteDistrict")
	beego.Router("/game/:id:int/add", &DistrictAddController{})

	sc := new(ServerController)
	beego.Router("/game/:id1:int/:id2:int", sc)
	beego.Router("/game/:id1:int/:id2:int/:id3:int", sc, "*:OpServer")
	beego.Router("/game/:id1:int/all", sc, "*:AllOpServer")
	beego.Router("/game/:id1:int/:id2:int/:id3:int/del", sc, "get:DeleteServer")
	beego.Router("/game/:id1:int/:id2:int/:id3:int/edit", sc, "get:ShowServer;post:UpdateServer")
	beego.Router("/game/:id1:int/:id2:int/add", &ServerAddController{})
}
