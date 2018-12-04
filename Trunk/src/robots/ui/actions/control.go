package actions

import "github.com/lunny/tango"
import "robots/controller"
import "os"
import "time"

type Control struct {
	RenderBase
	tango.Json
}

func (c *Control) Post() interface{} {
	type jsondata struct {
		Cmd  string `json:"cmd"`
		Args string `json:"args"`
	}
	var data jsondata
	err := c.DecodeJSON(&data)
	if err != nil {
		return map[string]interface{}{
			"err:":   err.Error(),
			"status": 500,
		}
	}

	switch data.Cmd {
	case "startall":
		controller.CreateRobot()
		controller.StartAll()
	case "stopall":
		controller.Shutdown()
	case "shutdown":
		controller.Shutdown()
		go func() {
			time.Sleep(time.Second * 2)
			os.Exit(0)
		}()
	case "queryinfo":
		ret, errs := controller.QueryInfo()
		return map[string]interface{}{
			"status":  200,
			"data":    ret,
			"errinfo": errs,
			//"sysinfo": controller.GetSysInfo(),
		}
	case "begin_move":
		controller.CommandBeginMove()
	case "stop_move":
		controller.CommandStopMove()
	case "begin_chat":
		controller.CommandBeginChat()
	case "stop_chat":
		controller.CommandStopChat()
	case "switch_scene":
		controller.CommandSwitchScene(data.Args)
	case "begin_attack":
		controller.CommandBeginAttack()
	case "stop_attack":
		controller.CommandStopAttack()
	case "moveto":
		controller.CommandMoveTo(data.Args)
	case "recover":
		controller.CommandRecover()
	case "add_exp":
		controller.CommandAddExp()
	case "buy":
		controller.CommandBuy()
	case "refresh_mall":
		controller.CommandRefreshMall()
	case "rich":
		controller.CommandRich()
	case "draw_prize":
		controller.CommandDrawPrize()
	case "all_open":
		controller.CommandAllOpen()
	case "enter_clone_scene":
		controller.CommandEnterCloneScene()
	case "quit_clone_scene":
		controller.CommandQuitCloneScene()
	case "join_guild":
		controller.CommandJoinGuild(data.Args)
	case "quit_guild":
		controller.CommandQuitGuild()
	case "power_up":
		controller.CommandPowerUp()
	case "power_down":
		controller.CommandPowerDown()
	case "send_custom":
		controller.CommandSendCustom(data.Args)
	case "send_gm":
		controller.CommandSendGM(data.Args)
	case "enter_15v15":
		controller.CommandEnter15V15()
	case "quit_15v15":
		controller.CommandQuit15V15()
	case "enter_clone_scene_details":
		controller.CommandEnterCloneSceneDetails(data.Args)
	case "get_clone_sceneid_list":
		list, errs := controller.CommandGetCloneSceneIdList()
		return map[string]interface{}{
			"status":  200,
			"data":    list,
			"errinfo": errs,
			//"sysinfo": controller.GetSysInfo(),
		}
	case "pk":
		controller.CommandPK()
	case "enter_multi_scene":
		controller.CommandEnterMultiScene(data.Args)
	case "createteam":
		controller.CommandCreateTeam()
	case "scene_move":
		controller.CommandSceneMove(data.Args)
	case "stop_scene_move":
		controller.CommandStopSceneMove()
	case "submit_quest":
		controller.CommandSubmitQuest()
	case "accept_quest":
		controller.CommandAcceptQuest()
	case "ninja_arena":
		controller.CommandNinjaArena()
	case "stop_ninja_arena":
		controller.CommandStopNinjaArena()
	case "group_pk":
		controller.CommandGroupPKMsg()
	case "team_pk":
		controller.CommandTeamPK()
	default:
		return map[string]interface{}{
			"status": 404,
		}
	}
	return map[string]interface{}{
		"status": 200,
	}
}
