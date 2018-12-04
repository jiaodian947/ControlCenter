package routers

import (
	"speedy/controllers"

	"github.com/astaxie/beego"
)

func init() {
	User := new(controllers.User)
	Mail := new(controllers.Mail)
	Account := new(controllers.Account)
	Activity := new(controllers.Activity)
	Player := new(controllers.Player)

	Prop := new(controllers.Prop)
	View := new(controllers.View)
	Server := new(controllers.Server)
	Logs := new(controllers.Logs)
	Statistics := new(controllers.Statistics)

	beego.Router("/", new(controllers.MainController))
	beego.Router("/api/user/login", new(controllers.Login))
	beego.Router("/api/user/logout", new(controllers.Logout))
	//user
	beego.Router("/api/user/info", User, "get:GetUserInfo")
	beego.Router("/api/user/viewsPower", User, "get:GetViewsPower")
	beego.Router("/api/user/serversPower", User, "get:GetServersPower")
	beego.Router("/api/user/allUsersInfo", User, "get:GetAllUsersInfo")
	beego.Router("/api/user/changeServerPower", User, "post:ChangeServerPower")
	beego.Router("/api/user/changeViewPower", User, "post:ChangeViewPower")
	beego.Router("/api/user/whiteList", User, "get:GetWhiteList")
	beego.Router("/api/user/deleteWhiteListItem", User, "get:DeleteWhiteListItem")
	beego.Router("/api/user/addWhiteListItem", User, "post:AddWhiteListItem")

	//mail
	beego.Router("/api/mail/sendMailList", Mail, "get:GetSendMailList")
	beego.Router("/api/mail/deleteMailLog", Mail, "get:DeleteSendMailLog")
	beego.Router("/api/mail/sendMail", Mail, "post:SendMail")
	//account
	beego.Router("/api/account/sendAccountList", Account, "get:GetSendAccountList")
	beego.Router("/api/account/deleteAccountLog", Account, "get:DeleteSendAccountLog")
	beego.Router("/api/account/ban", Account, "post:Ban")
	//activity
	beego.Router("/api/activity/switchList", Activity, "get:GetSwitchList")
	beego.Router("/api/activity/changeSwitchStatus", Activity, "post:ChangeSwitchStatus")
	beego.Router("/api/activity/addActivityAnnouncement", Activity, "post:AddActivityAnnouncement")
	beego.Router("/api/activity/deleteActivityAnnouncement", Activity, "get:DeleteActivityAnnouncement")

	beego.Router("/api/activity/getActivityAnnouncementList", Activity, "get:GetActivityAnnouncementList")
	beego.Router("/api/activity/addNotice", Activity, "post:AddNotice")
	beego.Router("/api/activity/deleteNotice", Activity, "get:DeleteNotice")
	beego.Router("/api/activity/getNoticeList", Activity, "get:GetNoticeList")
	beego.Router("/api/activity/changeQuestion", Activity, "post:ChangeQuestion")

	// Player
	beego.Router("/api/player/playerInfo", Player, "post:GetPlayerInfo")

	//props
	beego.Router("/api/props/getPropsList", Prop, "get:GetPropsList")
	//views
	beego.Router("/api/view/viewList", View, "get:GetViewList")
	beego.Router("/api/view/viewDetails", View, "get:GetViewByViewId")
	beego.Router("/api/view/addView", View, "post:AddView")

	//servers
	beego.Router("/api/server/serverList", Server, "get:GetServerList")
	beego.Router("/api/server/queryDomainInfo", Server, "post:GetDomainInfo")
	beego.Router("/api/server/getDomainRecordName", Server, "post:GetDomainRecordName")
	//logs
	beego.Router("/api/logs/getGameLogs", Logs, "post:GetGameLogs")
	beego.Router("/api/logs/getGuildFunctionLogs", Logs, "post:GetGuildFunctionLogs")
	beego.Router("/api/logs/getGuildGameLogs", Logs, "post:GetGuildGameLogs")

	//statistics
	beego.Router("/api/stat/getGameStatistics", Statistics, "post:GetGameStatistics")
	beego.Router("/api/stat/getCapitalStatistics", Statistics, "post:GetCapitalStatistics")
	beego.Router("/api/stat/getCapitalStatisticsDetails", Statistics, "post:GetCapitalStatisticsDetails")
	beego.Router("/api/stat/getQuestionNaireStatistics", Statistics, "get:GetQuestionNaireStatistics")
	beego.Router("/api/stat/getQuestionNaireStatisticsDetails", Statistics, "post:GetQuestionNaireStatisticsDetails")
	beego.Router("/api/stat/getQuestionNaireNotMustStatisticsDetails", Statistics, "post:GetQuestionNaireNotMustStatisticsDetails")

}
