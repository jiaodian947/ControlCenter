package controller

import (
	"encoding/json"
	"robots/robot"
)

var (
	Ctx *RobotCtl
)

type RobotInfos struct {
	Started   bool `json:"started"`
	Total     int  `json:"total"`
	Connected int  `json:"connected"`
	Ready     int  `json:"ready"`
	Errors    int  `json:"errors"`
}

type Error struct {
	Index   int    `json:"index"`
	Account string `json:"account"`
	Role    string `json:"role"`
	Err     string `json:"err"`
}

type Errors struct {
	Error []Error `json:"error"`
}

func QueryInfo() (info, errinfo string) {
	errs := Errors{make([]Error, 0, 64)}
	ri := RobotInfos{}
	ri.Started = Ctx.Started
	ri.Total = Ctx.count
	for _, r := range Ctx.robots {
		if r == nil {
			continue
		}
		s := r.State()
		if s >= robot.ROBOT_STATE_CONNECTED {
			ri.Connected++
		}

		switch s {
		case robot.ROBOT_STATE_READY:
			ri.Ready++
		case robot.ROBOT_STATE_ERROR:
			fallthrough
		case robot.ROBOT_STATE_FAILED:
			ri.Errors++
			e := Error{}
			e.Account = r.GetRobot().Account
			e.Role = r.GetRobot().Name
			e.Index = r.GetRobot().Index
			e.Err = r.GetRobot().Err.Error()
			errs.Error = append(errs.Error, e)
		}
	}
	b, _ := json.Marshal(&ri)
	info = string(b)
	b1, _ := json.Marshal(&errs)
	errinfo = string(b1)
	return
}

func SetServerInfo(address string, port int, serverid string) {
	Ctx.serverip = address
	Ctx.port = port
	Ctx.serverid = serverid
}

func SetAccInfo(accprefix, passwd, nameprefix string, accstart int, count int) {
	Ctx.accprefix = accprefix
	Ctx.passwd = passwd
	Ctx.nameprefix = nameprefix
	Ctx.accstart = accstart
	Ctx.count = count
}

func CreateRobot() {
	Ctx.CreateAll()
}

func StartAll() {
	Ctx.StartAll()
}

func Shutdown() {
	Ctx.DestroyAll()
}

func CommandBeginMove() {
	Ctx.BeginMove()
}

func CommandStopMove() {
	Ctx.StopMove()
}

func CommandBeginChat() {
	Ctx.BeginChat()
}

func CommandStopChat() {
	Ctx.StopChat()
}

func CommandSwitchScene(scene string) {
	Ctx.SwitchScene(scene)
}

func CommandBeginAttack() {
	Ctx.BeginAttack()
}

func CommandStopAttack() {
	Ctx.StopAttack()
}

func CommandMoveTo(dest string) {
	Ctx.MoveTo(dest)
}

func CommandRecover() {
	Ctx.Recover()
}

func CommandAddExp() {
	Ctx.AddExp()
}

func CommandBuy() {
	Ctx.Buy()
}

func CommandRefreshMall() {
	Ctx.RefreshMall()
}

func CommandRich() {
	Ctx.Rich()
}

func CommandDrawPrize() {
	Ctx.DrawPrize()
}

func CommandAllOpen() {
	Ctx.AllOpen()
}

func CommandEnterCloneScene() {
	Ctx.EnterCloneScene()
}

func CommandQuitCloneScene() {
	Ctx.QuitCloneScene()
}

func CommandJoinGuild(guild string) {
	Ctx.JoinGuild(guild)
}

func CommandQuitGuild() {
	Ctx.QuitGuild()
}

func CommandPowerUp() {
	Ctx.PowerUp()
}

func CommandPowerDown() {
	Ctx.PowerDown()
}

func CommandSendCustom(custom string) {
	Ctx.SendCustom(custom)
}

func CommandSendGM(gm string) {
	Ctx.SendGM(gm)
}

func CommandEnter15V15() {
	Ctx.Enter15V15()
}

func CommandQuit15V15() {
	Ctx.Quit15V15()
}
func CommandEnterCloneSceneDetails(sceneid string) {
	Ctx.EnterCloneSceneDetails(sceneid)
}
func CommandPK() {
	Ctx.PK()
}
func CommandGetCloneSceneIdList() (list []int64, error string) {
	return Ctx.GetCloneSceneIdList()
}
func CommandEnterMultiScene(sceneid string) {
	Ctx.EnterMultiScene(sceneid)
}
func CommandCreateTeam() {
	Ctx.CreateTeam()
}
func CommandSceneMove(custom string) {
	Ctx.SceneMove(custom)
}

func CommandStopSceneMove() {
	Ctx.StopSceneMove()
}
func CommandAcceptQuest() {
	Ctx.AcceptQuest()
}
func CommandSubmitQuest() {
	Ctx.SubmitQuest()
}
func CommandNinjaArena() {
	Ctx.NinjaArena()
}
func CommandStopNinjaArena() {
	Ctx.StopNinjaArena()
}

func CommandGroupPKMsg() {
	Ctx.GroupPK()
}

func CommandTeamPK() {
	Ctx.TeamPK()
}
