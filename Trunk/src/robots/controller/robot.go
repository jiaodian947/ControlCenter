package controller

import (
	"fmt"
	"robots/robot"
	"runtime"
	"runtime/debug"
	"strings"
)

type Robot interface {
	GetAttr(attr string) interface{}
	GetRobot() *robot.Robot
	Connect(addr string, port int, serverid string)
	Destroy()
	State() int
	BeginMove()
	StopMove()
	BeginChat()
	StopChat()
	SwitchScene(scene string)
	BeginAttack()
	StopAttack()
	MoveTo(dest string)
	Recover()
	AddExp()
	Buy()
	RefreshMall()
	BecomeRich()
	DrawPrize()
	AllOpen()
	EnterCloneScene()
	QuitCloneScene()
	JoinGuild(guild string)
	QuitGuild()
	PowerUp()
	PowerDown()
	SendCustomMsg(custom string)
	SendGM(gm string)
	Enter15V15()
	Quit15V15()
	EnterCloneSceneDetails(sceneid string)
	GetCloneSceneIdList() (list []int64, error string)
	PK(rguid int64)
	EnterMultiScene(sceneid string)
	CreateTeam()
	AutoMatchTeam()
	SceneMove(custom string)
	StopSceneMove()
	SubmitQuest()
	AcceptQuest()
	NinjaArena()
	StopNinjaArena()
	GroupPK()
	TeamPK()
}

type CreateFunc func(acc, pwd, name string, index int) Robot
type RobotCtl struct {
	robots     []Robot
	count      int
	serverip   string
	port       int
	serverid   string
	accprefix  string
	accstart   int
	passwd     string
	nameprefix string
	createfunc CreateFunc
	Started    bool
}

func New(f CreateFunc) {
	Ctx = &RobotCtl{}
	Ctx.robots = make([]Robot, 0, 2048)
	Ctx.createfunc = f
}

func (c *RobotCtl) CreateAll() {
	if len(c.robots) != 0 {
		return
	}

	for i := 0; i < c.count; i++ {
		name := fmt.Sprintf("%s_%s_%03d", c.nameprefix, c.serverid, c.accstart+i)
		if strings.Index(name, "64") != -1 {
			name = strings.Replace(name, "64", "L4", -1)
		}
		if strings.Index(name, "96") != -1 {
			name = strings.Replace(name, "96", "J6", -1)
		}
		acc := fmt.Sprintf("%s%d", c.accprefix, c.accstart+i)
		r := c.createfunc(acc, c.passwd, name, i)
		if r == nil {
			panic("create failed")
		}
		c.robots = append(c.robots, r)
	}
}

func (c *RobotCtl) StartAll() {
	for _, r := range c.robots {
		if r != nil {
			r.Connect(c.serverip, c.port, c.serverid)
		}
	}
	c.Started = true
}

func (c *RobotCtl) DestroyAll() {
	for _, r := range c.robots {
		if r != nil {
			r.Destroy()
		}
	}
	c.robots = c.robots[:0]
	c.Started = false
	runtime.GC()
	debug.FreeOSMemory()
}

func (c *RobotCtl) BeginMove() {
	for _, r := range c.robots {
		if r != nil && r.State() == robot.ROBOT_STATE_READY {
			r.BeginMove()
		}
	}
}

func (c *RobotCtl) StopMove() {
	for _, r := range c.robots {
		if r != nil {
			r.StopMove()
		}
	}
}

func (c *RobotCtl) BeginChat() {
	for _, r := range c.robots {
		if r != nil && r.State() == robot.ROBOT_STATE_READY {
			r.BeginChat()
		}
	}
}

func (c *RobotCtl) StopChat() {
	for _, r := range c.robots {
		if r != nil {
			r.StopChat()
		}
	}
}

func (c *RobotCtl) SwitchScene(scene string) {
	for _, r := range c.robots {
		if r != nil {
			r.SwitchScene(scene)
		}
	}
}

func (c *RobotCtl) BeginAttack() {
	for _, r := range c.robots {
		if r != nil && r.State() == robot.ROBOT_STATE_READY {
			r.BeginAttack()
		}
	}
}

func (c *RobotCtl) StopAttack() {
	for _, r := range c.robots {
		if r != nil {
			r.StopAttack()
		}
	}
}

func (c *RobotCtl) MoveTo(dest string) {
	for _, r := range c.robots {
		if r != nil {
			r.MoveTo(dest)
		}
	}
}

func (c *RobotCtl) Recover() {
	for _, r := range c.robots {
		if r != nil {
			r.Recover()
		}
	}
}

func (c *RobotCtl) AddExp() {
	for _, r := range c.robots {
		if r != nil {
			r.AddExp()
		}
	}
}

func (c *RobotCtl) Buy() {
	for _, r := range c.robots {
		if r != nil {
			r.Buy()
		}
	}
}

func (c *RobotCtl) RefreshMall() {
	for _, r := range c.robots {
		if r != nil {
			r.RefreshMall()
		}
	}
}

func (c *RobotCtl) Rich() {
	for _, r := range c.robots {
		if r != nil {
			r.BecomeRich()
		}
	}
}

func (c *RobotCtl) DrawPrize() {
	for _, r := range c.robots {
		if r != nil {
			r.DrawPrize()
		}
	}
}

func (c *RobotCtl) AllOpen() {
	for _, r := range c.robots {
		if r != nil {
			r.AllOpen()
		}
	}
}

func (c *RobotCtl) EnterCloneScene() {
	for _, r := range c.robots {
		if r != nil {
			r.EnterCloneScene()
		}
	}
}

func (c *RobotCtl) QuitCloneScene() {
	for _, r := range c.robots {
		if r != nil {
			r.QuitCloneScene()
		}
	}
}

func (c *RobotCtl) JoinGuild(guild string) {
	for _, r := range c.robots {
		if r != nil {
			r.JoinGuild(guild)
		}
	}
}

func (c *RobotCtl) QuitGuild() {
	for _, r := range c.robots {
		if r != nil {
			r.QuitGuild()
		}
	}
}

func (c *RobotCtl) PowerUp() {
	for _, r := range c.robots {
		if r != nil {
			r.PowerUp()
		}
	}
}

func (c *RobotCtl) PowerDown() {
	for _, r := range c.robots {
		if r != nil {
			r.PowerDown()
		}
	}
}

func (c *RobotCtl) SendCustom(args string) {
	for _, r := range c.robots {
		if r != nil {
			r.SendCustomMsg(args)
		}
	}
}

func (c *RobotCtl) SendGM(args string) {
	for _, r := range c.robots {
		if r != nil {
			r.SendGM(args)
		}
	}
}

func (c *RobotCtl) Enter15V15() {
	for _, r := range c.robots {
		if r != nil {
			r.Enter15V15()
		}
	}
}

func (c *RobotCtl) Quit15V15() {
	for _, r := range c.robots {
		if r != nil {
			r.Quit15V15()
		}
	}
}

func (c *RobotCtl) EnterCloneSceneDetails(args string) {
	for _, r := range c.robots {
		if r != nil {
			r.EnterCloneSceneDetails(args)
		}
	}
}
func (c *RobotCtl) PK() {

	for i := 0; i < len(c.robots); {
		c.robots[i].PK(c.robots[i+1].GetAttr("RGuid").(int64))
		i = i + 2
		if i > len(c.robots) {
			break
		}
	}

	// for _, r := range c.robots {
	// 	if r != nil {
	// 		rguid := r.GetAttr("RGuid")
	// 		fmt.Println(rguid)
	// 	}
	// }
}

func (c *RobotCtl) GetCloneSceneIdList() (list []int64, error string) {
	return c.robots[0].GetCloneSceneIdList()
}

func (c *RobotCtl) EnterMultiScene(args string) {
	for i := 0; i < len(c.robots); {
		c.robots[i].EnterMultiScene(args)
		i = i + 5
		if i > len(c.robots) {
			break
		}
	}
}
func (c *RobotCtl) CreateTeam() {
	for i := 0; i < len(c.robots); {
		c.robots[i].CreateTeam()
		for j := 1; j < 5; j++ {
			c.robots[i+j].AutoMatchTeam()
		}
		i = i + 5
		if i > len(c.robots) {
			break
		}
	}
}
func (c *RobotCtl) SceneMove(args string) {

	for _, r := range c.robots {
		if r != nil {
			r.SceneMove(args)
		}
	}
}
func (c *RobotCtl) StopSceneMove() {

	for _, r := range c.robots {
		if r != nil {
			r.StopSceneMove()
		}
	}
}
func (c *RobotCtl) SubmitQuest() {
	for _, r := range c.robots {
		if r != nil {
			r.SubmitQuest()
		}
	}
}
func (c *RobotCtl) AcceptQuest() {
	for _, r := range c.robots {
		if r != nil {
			r.AcceptQuest()
		}
	}
}
func (c *RobotCtl) NinjaArena() {
	for _, r := range c.robots {
		if r != nil {
			r.NinjaArena()
		}
	}
}
func (c *RobotCtl) StopNinjaArena() {
	for _, r := range c.robots {
		if r != nil {
			r.StopNinjaArena()
		}
	}
}

func (c *RobotCtl) GroupPK() {
	first := true
	for _, r := range c.robots {
		if r != nil {
			if first {
				r.SendGM("start40")
				first = false
			}
			r.GroupPK()
		}
	}
}

func (c *RobotCtl) TeamPK() {
	first := true
	for _, r := range c.robots {
		if r != nil {
			if first {
				r.SendGM("start5")
				first = false
			}
			r.TeamPK()
		}
	}
}
