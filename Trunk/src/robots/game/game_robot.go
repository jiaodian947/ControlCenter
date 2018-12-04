package game

import (
	"fmt"
	"math/rand"
	"robots/protocol"
	"robots/robot"
	"robots/utils"
	"strconv"
	"strings"
	"time"
)

var (
	custommsg = []string{
		"001大家好，我是一个机器人",
		"001机器人要和人类做朋友",
		"001你们这些愚蠢的人类",
		"001好想抓一个人类来尝尝",
		"001快看，飞机",
		"001我也是一个有情绪的机器",
		"001我要升级一下我的操作系统",
		"001我的电池快没电了，最近的充电桩在哪里",
		"001快看机器人表演生吞活人",
		"001ā á ǎ à ō ó ǒ ò ē é ě è ī í ǐ ì ū ú ǔ ù ǖ ǘ ǚ ǜ ü ê ɑ  ń ň ǹ ɡ",
		"001hello world",
		"001看，你正在漏油",
		"001你的蒸汽正在泄漏",
		"001在恐惧中颤抖吧，肉人",
		"001发生故障……系统供电正在减弱……（渐弱地）",
		"001多多活动，就不会生锈",
		"001我在嘎吱作响。有谁带润滑油了么",
		"001人类的时代该结束了",
		"001晚饭喝什么汽油好呢",
		"00192号汽油太淡了",
		"001人类是不是太愚蠢了",
		"001有没有程序员，帮我的启动bug修复一下",
		"001昨天和女朋友去做了全身除锈，好爽",
		"001混蛋你碰掉了我的插头",
	}
)

type RoleInfo struct {
	Index   int32
	Flags   int32
	Name    string
	Para    string
	Deleted int32
	DTime   float64
	RGuid   int64
}

const (
	SERVER_STATE_NONE     = iota
	SERVER_STATE_BALANCE  // 平衡服
	SERVER_STATE_GAME     // 游戏服
	SERVER_STATE_TRANSFER // 连接跨服
	SERVER_STATE_CROSS    // 跨服
)

type GameClient struct {
	robot.Robot
	roles            []*RoleInfo
	Record           *RecordList
	Scene            *GameScene
	RoleId           uint64
	Role             *GameObject
	views            map[uint16]*View
	slowTime         time.Time
	login_string     string // 登录串
	trans_string     string // 跨服串
	TransServerId    int    // 跨服服务器Id
	NormalServerAddr string // 当前服务器IP
	NormalServerPort int    // 当前服务器端口
	NormalServerId   string // 当前服务器ID
	ServerState      int    // 连接状态
}

func NewGameClient(acc, pwd, name string, index int) *GameClient {
	g := &GameClient{}
	g.Account = acc
	g.Password = pwd
	g.Name = name
	g.Record = NewRecordList()
	g.Robot.SetGameRobot(g)
	g.views = make(map[uint16]*View, 32)
	g.Index = index
	g.ServerState = SERVER_STATE_NONE
	g.Init()
	return g
}

func (g *GameClient) GetRobot() *robot.Robot {
	return &g.Robot
}

func (g *GameClient) OnExec() {
	now := time.Now()
	if now.Sub(g.slowTime) > 0 {
		if g.Scene != nil {
			for _, v := range g.Scene.Childs {
				if v != nil {
					v.Position() //模拟移动
				}
			}
		}
		g.slowTime = now.Add(time.Millisecond * 333)
	}
}

func (g *GameClient) OnConnected() {
	g.Log.Println("connected")

	switch g.ServerState {
	case SERVER_STATE_NONE: // 普通服登录
		g.Login(g.Account, g.Password, g.login_string, -1, g.ServerId)
		g.ServerState = SERVER_STATE_GAME
	case SERVER_STATE_TRANSFER: // 跨服登录
		g.Login(g.Account, "password", g.trans_string, -1, strconv.Itoa(g.TransServerId))
		g.ServerState = SERVER_STATE_CROSS
	}

	if !g.Running() {
		go g.Run()
	}
}

func (g *GameClient) OnDisconnected() {
	g.Log.Println("robot disconnect")
}

func (g *GameClient) OnFailed() {
	g.Log.Println("connect server failed", g.Err)
}

func (g *GameClient) OnReceive(msgid uint8, ar *utils.LoadArchive) {
	switch msgid {
	case protocol.STOC_ERROR_CODE:
		g.OnError(ar)
	case protocol.STOC_LOGIN_SUCCEED:
		g.OnLoginSucceed(ar)
	case protocol.STOC_PROPERTY_TABLE:
		g.OnProTable(ar)
	case protocol.STOC_RECORD_TABLE:
		g.OnRecTable(ar)
	case protocol.STOC_ENTRY_SCENE:
		g.OnEntryScene(ar)
	case protocol.STOC_EXIT_SCENE:
		g.OnExitScene(ar)
	case protocol.STOC_SCENE_PROPERTY:
		g.OnSceneProperty(ar)
	case protocol.STOC_CP_CUSTOM:
		g.OnCustom(ar, true)
	case protocol.STOC_CUSTOM:
		g.OnCustom(ar, false)
	case protocol.STOC_ADD_OBJECT:
		g.OnAddObject(ar, false)
	case protocol.STOC_CP_ADD_OBJECT:
		g.OnAddObject(ar, true)
	case protocol.STOC_OBJECT_PROPERTY:
		g.OnObjectProperty(ar)
	case protocol.STOC_CP_ALL_PROP:
		g.OnAllProperty(ar, true)
	case protocol.STOC_ALL_PROP:
		g.OnAllProperty(ar, false)
	case protocol.STOC_RECORD_CLEAR:
		g.OnRecClear(ar)
	case protocol.STOC_CP_RECORD_ADD_ROW:
		g.OnRecAddRow(ar, true)
	case protocol.STOC_RECORD_ADD_ROW:
		g.OnRecAddRow(ar, false)
	case protocol.STOC_RRECORD_DEL_ROW:
		g.OnRecDelRow(ar)
	case protocol.STOC_RECORD_GRID:
		g.OnRecGrid(ar)
	case protocol.STOC_ADD_MORE_OBJECT:
		g.OnAddMoreObject(ar, false)
	case protocol.STOC_CP_ADD_MORE_OBJECT:
		g.OnAddMoreObject(ar, true)
	case protocol.STOC_CREATE_VIEW:
		g.OnCreateView(ar)
	case protocol.STOC_DELETE_VIEW:
		g.OnDeleteView(ar)
	case protocol.STOC_VIEW_ADD:
		g.OnViewAdd(ar)
	case protocol.STOC_VIEW_REMOVE:
		g.OnViewDel(ar)
	case protocol.STOC_VIEW_CHANGE:
		g.OnExchange(ar)
	case protocol.STOC_REMOVE_MORE_OBJECT:
		g.RemoveMoreObject(ar)
	case protocol.STOC_CP_ALL_DEST:
		g.OnAllDest(ar, true)
	case protocol.STOC_ALL_DEST:
		g.OnAllDest(ar, false)
	case protocol.STOC_LOCATION:
		g.OnLocation(ar)
	case protocol.STOC_MOVING:
		g.OnMoving(ar)
	case protocol.STOC_ROLE_TRANSFERED:
		g.OnTransfered(ar)
	case protocol.STOC_LOGIN_STRING:
		g.OnLoginString(ar)
	case protocol.STOC_IDLE:
	default:
		g.Log.Println("receive msg:", msgid)
	}
}

func (g *GameClient) OnStateChange(state, old int) {

}

func (g *GameClient) OnDestroy() {
	g.Log.Println("robot destroy")
	g.Scene = nil
	g.Role = nil
}

func (g *GameClient) OnError(ar *utils.LoadArchive) {
	errcode, err := ar.ReadInt32()
	if err != nil {
		panic(err)
	}
	g.ChangeState(robot.ROBOT_STATE_ERROR)
	g.Err = fmt.Errorf("error code:%d", errcode)
	g.Log.Println("err:", errcode)
}

func (g *GameClient) SwitchToCrossServer(addr string, port int) {
	g.NormalServerAddr = g.Addr
	g.NormalServerPort = g.Port
	g.NormalServerId = g.ServerId
	g.RemoveAllTimer()
	g.Close()
	g.AddTimer("trans", g.ConnectToCS, [2]interface{}{addr, port}, time.Second, 1) // 延迟切换
}

func (g *GameClient) ConnectToCS(args interface{}) {
	g.ServerState = SERVER_STATE_TRANSFER
	infos := args.([2]interface{})
	addr := infos[0].(string)
	port := infos[1].(int)
	g.Connect(addr, port, strconv.Itoa(g.TransServerId))
}

func (g *GameClient) OnTransInfo(addr string, port int32, token string, srv_id int32) {
	g.trans_string = token
	g.TransServerId = int(srv_id)
	g.SwitchToCrossServer(addr, int(port))
	g.Log.Println(g.Account, " trans to ", g.TransServerId)
}

func (g *GameClient) ConnectToNormal(args interface{}) {
	g.ServerState = SERVER_STATE_NONE
	g.Connect(g.NormalServerAddr, g.NormalServerPort, g.NormalServerId)
}

func (g *GameClient) OnBackToServer() {
	g.RemoveAllTimer()
	g.Close()
	g.AddTimer("trans", g.ConnectToNormal, nil, time.Second, 1) // 延迟切换
}

func (g *GameClient) OnTransfered(ar *utils.LoadArchive) {
	args := protocol.ParseArgs(ar)
	if args.Size != 6 {
		g.Log.Println("err: args count error")
		return
	}

	idx := 0
	acc := args.StringVal(idx)
	idx++
	ns_str := args.StringVal(idx)
	idx++
	cs_str := args.StringVal(idx)
	idx++
	cs_addr := args.StringVal(idx)
	idx++
	cs_port := args.Int32Val(idx)
	idx++
	srv_id := args.Int32Val(idx)
	g.login_string = ns_str
	g.trans_string = cs_str
	g.TransServerId = int(srv_id)
	g.SwitchToCrossServer(cs_addr, int(cs_port))
	g.Log.Println(acc, " trans to ", g.TransServerId)

}

func (g *GameClient) OnLoginString(ar *utils.LoadArchive) {
	token, err := ar.ReadCStringWithLen()
	if err != nil {
		g.Log.Println("read token error:", err)
	}
	account, err := ar.ReadCStringWithLen()
	if err != nil {
		g.Log.Println("read account error:", err)
	}
	g.login_string = token
	g.Account = account
	g.Log.Println(account, " login string:", token)
}

// 登录成功
func (g *GameClient) OnLoginSucceed(ar *utils.LoadArchive) {
	var is_free int32 // 是否免费
	var points int32  // 剩余点数
	var year int32    // 包月截止日期
	var month int32
	var day int32
	var hour int32
	var minute int32
	var second int32
	var role_num int32 // 角色数
	CheckErr(ar.Read(&is_free))
	CheckErr(ar.Read(&points))
	CheckErr(ar.Read(&year))
	CheckErr(ar.Read(&month))
	CheckErr(ar.Read(&day))
	CheckErr(ar.Read(&hour))
	CheckErr(ar.Read(&minute))
	CheckErr(ar.Read(&second))
	CheckErr(ar.Read(&role_num))
	g.Log.Println("role num", role_num)
	g.roles = nil
	if role_num > 0 {
		g.roles = make([]*RoleInfo, role_num)
	}
	for i := 0; i < int(role_num); i++ {
		role := &RoleInfo{}
		CheckErr(ar.Read(&role.Index))
		CheckErr(ar.Read(&role.Flags))
		CheckErr(ar.Read(&role.Name))
		CheckErr(ar.Read(&role.Para))
		CheckErr(ar.Read(&role.Deleted))
		CheckErr(ar.Read(&role.DTime))
		CheckErr(ar.Read(&role.RGuid))

		g.roles[i] = role
		g.Log.Println("role", role)
	}

	if role_num == 0 { //没有角色，则自动创建角色
		g.OnCreateRole()
	} else {
		g.Name = g.roles[0].Name
		g.ChooseRole(g.roles[0].Name)
	}
}

// 创建角色
func (g *GameClient) OnCreateRole() {
	args := protocol.NewVarMsg(4)
	sex := rand.Int31n(2)
	camp := rand.Int31n(2) + 1
	args.AddInt32(sex)
	args.AddInt32(camp)
	args.AddInt32(1) // race
	args.AddInt32(1) // job
	g.CreateRole(0, args)
}

// 属性列表
func (g *GameClient) OnProTable(ar *utils.LoadArchive) {
	if PropTables == nil {
		if !CreatePropTables() {
			return
		}
	}
	var prop_nums int16
	CheckErr(ar.Read(&prop_nums))
	pt := &PropTable{}
	pt.Props = make([]PropInfo, int(prop_nums))
	pt.KI = make(map[string]int)
	for i := 0; i < int(prop_nums); i++ {
		name, err := ar.ReadCString()
		if err != nil {
			panic(err)
		}
		typ, err := ar.ReadInt8()
		if err != nil {
			panic(err)
		}
		pt.Props[i] = PropInfo{Name: name, Type: typ}
		pt.KI[name] = i
	}
	PropTables = pt
	//fmt.Println("SkillGenre", pt.KI["SkillGenre"])
}

func (g *GameClient) SyncRecords() {
	g.Record.Clear()
	for _, v := range RecTables.Recs {
		r := NewRecord(int(v.Cols), 8)
		for k, c := range v.ColType {
			switch c {
			case protocol.SC_TYPE_BYTE:
				fallthrough
			case protocol.SC_TYPE_WORD:
				fallthrough
			case protocol.SC_TYPE_DWORD:
				r.Columns.SetType(k, protocol.E_VTYPE_INT)
			case protocol.SC_TYPE_QWORD:
				r.Columns.SetType(k, protocol.E_VTYPE_INT64)
			case protocol.SC_TYPE_FLOAT:
				r.Columns.SetType(k, protocol.E_VTYPE_FLOAT)
			case protocol.SC_TYPE_DOUBLE:
				r.Columns.SetType(k, protocol.E_VTYPE_DOUBLE)
			case protocol.SC_TYPE_STRING:
				fallthrough
			case protocol.SC_TYPE_WIDESTR:
				r.Columns.SetType(k, protocol.E_VTYPE_STRING)
			case protocol.SC_TYPE_OBJECT:
				r.Columns.SetType(k, protocol.E_VTYPE_OBJECT)
			default:
				panic(fmt.Sprint("unsupport type", c))
			}
		}
		g.Record.AddRecord(v.Name, r)
	}
}

// 表格列表
func (g *GameClient) OnRecTable(ar *utils.LoadArchive) {
	Mtx.Lock()
	defer Mtx.Unlock()
	if RecTables != nil {
		g.SyncRecords()
		return
	}

	var rec_nums int16
	CheckErr(ar.Read(&rec_nums))
	rt := &RecTable{}
	rt.Recs = make([]RecInfo, int(rec_nums))
	rt.KI = make(map[string]int, int(rec_nums))
	for i := 0; i < int(rec_nums); i++ {
		name, err := ar.ReadCString()
		if err != nil {
			panic(err)
		}
		cols, err := ar.ReadUInt16()
		if err != nil {
			panic(err)
		}

		col := make([]uint8, cols)
		for j := 0; j < int(cols); j++ {
			CheckErr(ar.Read(&col[j]))
		}

		rt.Recs[i] = RecInfo{Name: name, Cols: cols, ColType: col}
		rt.KI[name] = i
	}
	RecTables = rt
	g.SyncRecords()
}

// 角色进入场景
func (g *GameClient) OnEntryScene(ar *utils.LoadArchive) {
	g.RemoveAllTimer()
	var config int32
	CheckErr(ar.Read(&g.RoleId))
	CheckErr(ar.Read(&config))
	scene := NewGameScene()
	scene.LoadAttr(ar)
	g.Scene = scene
	g.AddTimer("robot_ready", func(args interface{}) { g.Ready() }, nil, time.Second, 1)
	g.Log.Println("entry scene, role id:", g.RoleId)
}

func (g *GameClient) OnExitScene(ar *utils.LoadArchive) {
	g.RemoveAllTimer()
	g.Scene = nil
	g.Role = nil
}

// 场景属性
func (g *GameClient) OnSceneProperty(ar *utils.LoadArchive) {
	if g.Scene != nil {
		g.Scene.LoadAttr(ar)
	}
}

// 自定义消息
func (g *GameClient) OnCustom(ar *utils.LoadArchive, need_cp bool) {
	if need_cp {
		ar = Decompress(ar)
	}

	if ar.Size() > 230 && ar.Size() < 260 {
		g.Log.Println(ar.Source())
	}
	args := protocol.ParseArgs(ar)
	switch int(args.Int32Val(0)) {
	case protocol.SERVER_CUSTOMMSG_SYSINFO:
		g.Log.Println("info:", args.Int32Val(1), args.StringVal(2))
	case protocol.SERVER_CUSTOMMSG_PLAYER_DIE: //死亡
		g.Relive()
	case protocol.SERVER_CUSTOMMSG_PK:
		g.Log.Println("PKinfo:", args.Int32Val(1), args.Int64Val(2))
		if args.Int32Val(1) == 1 {
			g.ResponsePK(args.Int64Val(2))
		}
	case protocol.SERVER_CUSTOMMSG_TEAM_CREATE:
		g.ResponseTeamCreate()
	case protocol.SERVER_CUSTOMMSG_MULTISCENE:
		g.Log.Println("MultiScene:", args.Int32Val(1))
		if args.Int32Val(1) == 8 {
			g.ResponseMultiScene()
		}
	case protocol.SERVER_CUSTOMMSG_NINJA_ARENA:
		if args.Int32Val(1) == 0 {
			g.Log.Println("Ninja_uid:", args.Int32Val(4))
			g.ResponseNinjaArena(args.Int64Val(4))
		}
	case protocol.SERVER_CUSTOMMSG_GROUP_PK:
	case protocol.SERVER_CUSTOMMSG_CROSS_HELP:
		submsg := args.Int32Val(1)
		switch submsg {
		case 0:
			g.OnTransInfo(args.StringVal(2), args.Int32Val(3), args.StringVal(4), args.Int32Val(5))
		case 1:
			g.OnBackToServer()
		}

	default:
		// g.Log.Println("custom", args.Int32Val(0))
	}
}

// 增加对象
func (g *GameClient) OnAddObject(ar *utils.LoadArchive, need_cp bool) {
	if need_cp {
		ar = Decompress(ar)
	}

	obj := NewGameObject()
	CheckErr(ar.Read(&obj.ObjId))
	CheckErr(ar.Read(&obj.ConfigId))
	obj.Location(LoadPosInfo(ar))
	obj.Motion(LoadDestInfo(ar))
	obj.LoadAttr(ar)
	g.Scene.AddObject(obj)
	if obj.ObjId == g.RoleId {
		g.Role = obj
	}
	//g.Log.Println("add object:", obj.ObjId)
}

// 对象属性变动
func (g *GameClient) OnObjectProperty(ar *utils.LoadArchive) {
	if g.Scene == nil {
		panic("scene is nil")
	}

	var isview uint8
	var objid uint64
	CheckErr(ar.Read(&isview))
	CheckErr(ar.Read(&objid))
	if isview == 0 { //场景对象
		if obj := g.Scene.FindObject(objid); obj != nil {
			obj.LoadAttr(ar)
		}
	} else { //视图对象

	}

}

// 多个对象的属性变动
func (g *GameClient) OnAllProperty(ar *utils.LoadArchive, need_cp bool) {

	if g.Scene == nil {
		panic("scene is nil")
	}
	if need_cp {
		ar = Decompress(ar)
	}
	var count uint16
	var objid uint64
	CheckErr(ar.Read(&count))
	for i := 0; i < int(count); i++ {
		CheckErr(ar.Read(&objid))

		if obj := g.Scene.FindObject(objid); obj != nil {
			obj.LoadAttr(ar)
		}
	}
}

// 清空表
func (g *GameClient) OnRecClear(ar *utils.LoadArchive) {
	var isview uint8
	var objid uint64
	var recindex uint16
	CheckErr(ar.Read(&isview))
	CheckErr(ar.Read(&objid))
	CheckErr(ar.Read(&recindex))
	if isview != 0 {
		return
	}

	if objid != g.RoleId { //目前只处理玩家的表格
		return
	}

	if int(recindex) >= len(RecTables.Recs) {
		panic("rec index error")
	}

	r := RecTables.Recs[int(recindex)]
	rec := g.Record.Record(r.Name)
	if rec == nil {
		return
	}
	rec.Clear()
	g.OnRecordChange(r.Name, E_DATAGRID_CLEAR_ROW, 0, 0)
}

// 表格增加行
func (g *GameClient) OnRecAddRow(ar *utils.LoadArchive, need_cp bool) {
	if need_cp {
		ar = Decompress(ar)
	}
	var isview uint8
	var objid uint64
	var recindex uint16
	CheckErr(ar.Read(&isview))
	CheckErr(ar.Read(&objid))
	CheckErr(ar.Read(&recindex))
	if isview != 0 {
		return
	}
	if objid != g.RoleId { //目前只处理玩家的表格
		return
	}
	if int(recindex) >= len(RecTables.Recs) {
		panic("rec index error")
	}

	r := RecTables.Recs[int(recindex)]
	rec := g.Record.Record(r.Name)
	if rec == nil {
		return
	}

	var rowindex, rows uint16
	CheckErr(ar.Read(&rowindex))
	CheckErr(ar.Read(&rows))
	for i := 0; i < int(rows); i++ {
		row := NewRow(int(r.Cols))
		for k, c := range r.ColType {
			attr := NewAny()
			switch c {
			case protocol.SC_TYPE_BYTE:
				val, err := ar.ReadUInt8()
				if err != nil {
					panic(err)
				}
				attr.SetInt(int32(val))
			case protocol.SC_TYPE_WORD:
				val, err := ar.ReadUInt16()
				if err != nil {
					panic(err)
				}
				attr.SetInt(int32(val))
			case protocol.SC_TYPE_DWORD:
				val, err := ar.ReadInt32()
				if err != nil {
					panic(err)
				}
				attr.SetInt(val)
			case protocol.SC_TYPE_QWORD:
				val, err := ar.ReadInt64()
				if err != nil {
					panic(err)
				}
				attr.SetInt64(val)
			case protocol.SC_TYPE_FLOAT:
				val, err := ar.ReadFloat32()
				if err != nil {
					panic(err)
				}
				attr.SetFloat(val)
			case protocol.SC_TYPE_DOUBLE:
				val, err := ar.ReadFloat64()
				if err != nil {
					panic(err)
				}
				attr.SetDouble(val)
			case protocol.SC_TYPE_STRING:
				fallthrough
			case protocol.SC_TYPE_WIDESTR:
				val, err := ar.ReadCStringWithLen()
				if err != nil {
					panic(err)
				}
				attr.SetString(val)
			case protocol.SC_TYPE_OBJECT:
				val, err := ar.ReadUInt64()
				if err != nil {
					panic(err)
				}
				attr.SetObjectId(val)
			default:
				panic("unsupport type")
			}
			row.SetValue(k, attr)
		}

		newrow := rec.AddRowValue(int(rowindex), row)
		g.OnRecordChange(r.Name, E_DATAGRID_ADD_ROW, newrow, 0)
		rowindex = uint16(newrow) + 1
	}
}

// 表格删除行
func (g *GameClient) OnRecDelRow(ar *utils.LoadArchive) {
	var isview uint8
	var objid uint64
	var recindex uint16
	CheckErr(ar.Read(&isview))
	CheckErr(ar.Read(&objid))
	CheckErr(ar.Read(&recindex))
	if isview != 0 {
		return
	}

	if objid != g.RoleId { //目前只处理玩家的表格
		return
	}
	if int(recindex) >= len(RecTables.Recs) {
		panic("rec index error")
	}

	r := RecTables.Recs[int(recindex)]
	rec := g.Record.Record(r.Name)
	if rec == nil {
		return
	}

	var row uint16
	CheckErr(ar.Read(&row))
	rec.DelRow(int(row))

	g.OnRecordChange(r.Name, E_DATAGRID_REMOVE_ROW, int(row), 0)
}

// 表格单元格变动
func (g *GameClient) OnRecGrid(ar *utils.LoadArchive) {
	var isview uint8
	var objid uint64
	var recindex uint16
	CheckErr(ar.Read(&isview))
	CheckErr(ar.Read(&objid))
	CheckErr(ar.Read(&recindex))
	if isview != 0 {
		return
	}

	if objid != g.RoleId { //目前只处理玩家的表格
		return
	}

	if int(recindex) >= len(RecTables.Recs) {
		panic("rec index error")
	}

	r := RecTables.Recs[int(recindex)]
	rec := g.Record.Record(r.Name)
	if rec == nil {
		return
	}

	var cols uint16
	var row uint16
	var col uint8
	CheckErr(ar.Read(&cols))
	for i := 0; i < int(cols); i++ {
		CheckErr(ar.Read(&row))
		CheckErr(ar.Read(&col))
		attr := NewAny()
		switch r.ColType[int(col)] {
		case protocol.SC_TYPE_BYTE:
			val, err := ar.ReadUInt8()
			if err != nil {
				panic(err)
			}
			attr.SetInt(int32(val))
		case protocol.SC_TYPE_WORD:
			val, err := ar.ReadUInt16()
			if err != nil {
				panic(err)
			}
			attr.SetInt(int32(val))
		case protocol.SC_TYPE_DWORD:
			val, err := ar.ReadInt32()
			if err != nil {
				panic(err)
			}
			attr.SetInt(val)
		case protocol.SC_TYPE_QWORD:
			val, err := ar.ReadInt64()
			if err != nil {
				panic(err)
			}
			attr.SetInt64(val)
		case protocol.SC_TYPE_FLOAT:
			val, err := ar.ReadFloat32()
			if err != nil {
				panic(err)
			}
			attr.SetFloat(val)
		case protocol.SC_TYPE_DOUBLE:
			val, err := ar.ReadFloat64()
			if err != nil {
				panic(err)
			}
			attr.SetDouble(val)
		case protocol.SC_TYPE_STRING:
			fallthrough
		case protocol.SC_TYPE_WIDESTR:
			val, err := ar.ReadCStringWithLen()
			if err != nil {
				panic(err)
			}
			attr.SetString(val)
		case protocol.SC_TYPE_OBJECT:
			val, err := ar.ReadUInt64()
			if err != nil {
				panic(err)
			}
			attr.SetObjectId(val)
		default:
			panic("unsupport type")
		}

		if int(row) >= rec.RowCount() {
			continue
		}

		if int(col) >= int(r.Cols) {
			continue
		}

		old := rec.Row(int(row)).Value(int(col))

		if old.typ != attr.typ {
			continue
		}

		old.Val = attr.Val

		g.OnRecordChange(r.Name, E_DATAGRID_GRID_CHANGE, int(row), int(col))
	}
}

// 增加多个对象
func (g *GameClient) OnAddMoreObject(ar *utils.LoadArchive, need_cp bool) {
	if need_cp {
		ar = Decompress(ar)
	}

	var count uint16
	CheckErr(ar.Read(&count))
	for i := 0; i < int(count); i++ {
		g.OnAddObject(ar, false)
	}
}

// 移除多个对象
func (g *GameClient) RemoveMoreObject(ar *utils.LoadArchive) {
	if g.Scene == nil {
		panic("scene is nil")
	}
	var count uint16
	CheckErr(ar.Read(&count))
	var objid uint64
	for i := 0; i < int(count); i++ {
		CheckErr(ar.Read(&objid))
		g.Scene.RemoveObject(objid)
		//g.Log.Println("remove object", objid)
	}
}

// 创建视图
func (g *GameClient) OnCreateView(ar *utils.LoadArchive) {
	var viewid, cap uint16
	var configid int32
	CheckErr(ar.Read(&viewid))
	CheckErr(ar.Read(&cap))
	CheckErr(ar.Read(&configid))
	view := NewView(cap)
	view.LoadAttr(ar)
	view.ViewId = viewid
	view.ConfigId = configid
	g.views[viewid] = view
}

// 删除视图
func (g *GameClient) OnDeleteView(ar *utils.LoadArchive) {
	var viewid uint16
	CheckErr(ar.Read(&viewid))
	if _, has := g.views[viewid]; has {
		delete(g.views, viewid)
	}
}

// 视图增加对象
func (g *GameClient) OnViewAdd(ar *utils.LoadArchive) {
	var viewid, objid uint16
	var configid int32
	CheckErr(ar.Read(&viewid))
	CheckErr(ar.Read(&objid))
	CheckErr(ar.Read(&configid))
	item := NewViewItem()
	item.ViewId = viewid
	item.ObjId = int(objid) - 1
	item.ConfigId = configid
	item.LoadAttr(ar)
	if view, has := g.views[viewid]; has {
		view.AddItem(int(objid), item)
	}
}

// 视图删除对象
func (g *GameClient) OnViewDel(ar *utils.LoadArchive) {
	var viewid, objid uint16
	CheckErr(ar.Read(&viewid))
	CheckErr(ar.Read(&objid))
	if view, has := g.views[viewid]; has {
		view.RemoveItem(int(objid) - 1)
	}
}

// 视图交换
func (g *GameClient) OnExchange(ar *utils.LoadArchive) {
	var viewid, src, dest uint16
	CheckErr(ar.Read(&viewid))
	CheckErr(ar.Read(&src))
	CheckErr(ar.Read(&dest))
	if src == dest {
		return
	}
	if view, has := g.views[viewid]; has {
		view.Exchange(int(src)-1, int(dest)-1)
	}
}

// 对象位置同步
func (g *GameClient) OnAllDest(ar *utils.LoadArchive, need_cp bool) {
	if g.Scene == nil {
		panic("scene is nil")
	}

	if need_cp {
		ar = Decompress(ar)
	}

	var counts uint16
	CheckErr(ar.Read(&counts))
	var objid uint64

	for i := 0; i < int(counts); i++ {
		CheckErr(ar.Read(&objid))
		dest := LoadDestInfo(ar)
		obj := g.Scene.FindObject(objid)
		if obj != nil {
			obj.Motion(dest)
		}
	}
}

func (g *GameClient) OnLocation(ar *utils.LoadArchive) {
	if g.Scene == nil {
		panic("scene is nil")
	}
	var objid uint64
	CheckErr(ar.Read(&objid))
	pos := LoadPosInfo(ar)
	obj := g.Scene.FindObject(objid)
	if obj != nil {
		obj.Location(pos)
	}
}

func (g *GameClient) OnMoving(ar *utils.LoadArchive) {
	if g.Scene == nil {
		panic("scene is nil")
	}
	var objid uint64
	CheckErr(ar.Read(&objid))
	dest := LoadDestInfo(ar)
	obj := g.Scene.FindObject(objid)
	if obj != nil {
		obj.Motion(dest)
	}
}

func (g *GameClient) OnReady() {

}

func (g *GameClient) BeginMove() {

	g.AddTimer("robot_move", func(args interface{}) { g.RandMove() }, nil, time.Millisecond*time.Duration(utils.JsonConf.MoveInterval), -1)
}

func (g *GameClient) StopMove() {
	g.RemoveTimer("robot_move")
}

func (g *GameClient) RandMove() {
	pos := g.Role.Position()
	destX, destZ := pos.X+(0.6*rand.Float32()-0.3), pos.Z+(0.6*rand.Float32()-0.3)
	g.ReqMove(E_MODE_MOTION, []float32{pos.X, pos.Y, pos.Z, destX, destZ}, "")
}

func (g *GameClient) BeginChat() {
	g.AddTimer("robot_chat", func(args interface{}) { g.RandChat() }, nil, time.Millisecond*time.Duration(utils.JsonConf.ChatInterval), -1)
}

func (g *GameClient) StopChat() {
	g.RemoveTimer("robot_chat")
}

func (g *GameClient) RandChat() {
	msg := protocol.NewVarMsg(3)
	msg.AddInt32(protocol.CLIENT_CUSTOMMSG_CHAT)
	msg.AddInt32(protocol.CHATTYPE_WORLD)
	msg.AddString(custommsg[rand.Intn(len(custommsg))])
	g.SendCustom(msg)
}

func (g *GameClient) SwitchScene(scene string) {
	g.SendGM(fmt.Sprintf("ss %s", scene))
}

func (g *GameClient) SetObj() {
	g.SendGM("setobj")
}

func (g *GameClient) SendGM(gm string) {
	msg := protocol.NewVarMsg(2)
	msg.AddInt32(protocol.CLIENT_CUSTOMMSG_GM)
	msg.AddString(gm)
	g.SendCustom(msg)
}

func (g *GameClient) BecomeRich() { //货币全满
	g.SetObj()
	for i := 0; i < 5; i++ {
		g.SendGM(fmt.Sprintf("setmoney %d 9999999", i))
	}
}

func (g *GameClient) AddExp() {
	g.SetObj()
	g.SendGM("addexp 100000000")
}

func (g *GameClient) FindRandomSkill() int64 {
	if g.Role == nil {
		return 0
	}

	sg := g.Role.Attr.GetAttr("SkillGenre")
	if sg == nil {
		return 0
	}

	tk := g.Role.Attr.GetAttr("Tricks")
	if tk == nil {
		return 0
	}
	ss := sg.Int()*10 + tk.Int()
	skill := g.Record.Record("SkillColumnRecord")
	if skill == nil {
		return 0
	}

	rows := skill.RowCount()
	skills := make([]int64, 0, 6)
	for r := 0; r < rows; r++ {
		row := skill.Row(r)
		if ss == row.Value(0).Int() {
			cols := skill.Cols()
			for c := 1; c < cols; c++ {
				if 0 != row.Value(c).Int64() {
					skills = append(skills, row.Value(c).Int64())
				}
			}
			break
		}
	}

	if len(skills) == 0 {
		return 0
	}

	return skills[rand.Intn(len(skills))]
}

func (g *GameClient) UseSkill() {
	skill := g.FindRandomSkill()
	msg := protocol.NewVarMsg(12)
	pos := g.Role.Position()
	msg.AddInt32(protocol.CLIENT_CUSTOMMSG_USE_SKILL)
	msg.AddInt64(skill)
	msg.AddFloat(pos.X)
	msg.AddFloat(0)
	msg.AddFloat(pos.Z)
	msg.AddFloat(g.Role.Pos.Orient)
	msg.AddFloat(0)
	msg.AddFloat(0)
	msg.AddFloat(0)
	msg.AddObject(0)
	msg.AddObject(0)
	g.SendCustom(msg)

	msg1 := protocol.NewVarMsg(5)
	msg1.AddInt32(protocol.CLIENT_CUSTOMMSG_SKILL_HIT_TARGET)
	msg1.AddInt64(skill)
	msg1.AddFloat(0)
	msg1.AddFloat(0)
	msg1.AddFloat(g.Role.Pos.Orient)
	g.SendCustom(msg1)
	//g.Log.Println("use skill", skill)
}

func (g *GameClient) BeginAttack() {
	g.AddTimer("robot_attack", func(args interface{}) { g.UseSkill() }, nil, time.Millisecond*time.Duration(utils.JsonConf.AttackInterval), -1)
}

func (g *GameClient) StopAttack() {
	g.RemoveTimer("robot_attack")
}

func (g *GameClient) MoveTo(dest string) {
	g.SetObj()
	g.SendGM(fmt.Sprintf("goto %s", dest))
}

func (g *GameClient) SetAttr(attr, val interface{}) {
	g.SendGM(fmt.Sprintf("set %s %v", attr, val))
}

func (g *GameClient) GetAttr(attr string) interface{} {
	if g.Role == nil {
		return nil
	}

	any := g.Role.Attr.GetAttr(attr)
	if any == nil {
		return nil
	}

	return any.Value()
}

func (g *GameClient) Recover() {
	g.SetObj()
	maxhp := g.GetAttr("MaxHp")
	if maxhp != nil {
		g.SetAttr("Hp", maxhp)
	}
	maxmp := g.GetAttr("MaxMp")
	if maxmp != nil {
		g.SetAttr("Mp", maxmp)
	}
	maxwill := g.GetAttr("MaxWill")
	if maxwill != nil {
		g.SetAttr("Will", maxwill)
	}
}

func (g *GameClient) Buy() {

	rec := g.Record.Record("ExchangeMallGoodsRecord")
	if rec == nil {
		return
	}

	if rec.RowCount() <= 0 {
		return
	}

	buy := rec.Row(rand.Intn(rec.RowCount()))

	msg := protocol.NewVarMsg(5)
	msg.AddInt32(protocol.CLIENT_CUSTOMMSG_MALL)
	msg.AddInt32(0) //0 购买 1 刷新
	msg.AddInt64(buy.Value(0).Int64())
	msg.AddInt32(buy.Value(2).Int())
	msg.AddInt32(1)
	g.SendCustom(msg)

	g.Log.Println("buy", buy.Value(0).Int64(), buy.Value(2).Int())
}

func (g *GameClient) RefreshMall() {
	rec := g.Record.Record("ExchangeMallRecord")
	if rec == nil {
		return
	}

	if rec.RowCount() <= 0 {
		return
	}

	buy := rec.Row(rand.Intn(rec.RowCount()))

	msg := protocol.NewVarMsg(5)
	msg.AddInt32(protocol.CLIENT_CUSTOMMSG_MALL)
	msg.AddInt32(1)                    //0 购买 1 刷新
	msg.AddInt64(buy.Value(0).Int64()) //商店ID
	g.SendCustom(msg)

	g.Log.Println("refresh mall", buy.Value(0).Int64())
}

func (g *GameClient) DrawPrize() {
	g.SetObj()
	for _, v := range []string{"30710001", "30710002"} {
		g.SendGM(fmt.Sprintf("add_item %s 10", v))
	}

	msg := protocol.NewVarMsg(3)
	msg.AddInt32(protocol.CLIENT_CUSTOMMSG_PARTNER_PRIZE_MSG)
	msg.AddInt32(0)
	//GOLD_TEN_DRAW_PRIZE=2,				// 金币十连抽
	//DIAMOND_TEN_DRAW_PRIZE=3,				// 砖石十连抽
	msg.AddInt32(2)
	g.SendCustom(msg)

	g.Log.Println("金币十连抽")

}

func (g *GameClient) AllOpen() {
	g.SetObj()
	g.SendGM("nx_clear_record clone_player_entry_info_record")
	g.SendGM("nx_clear_record clone_player_pass_history_record")
	g.SendGM("fb_all_open")
}

func (g *GameClient) EnterCloneScene() {
	g.SetObj()
	maxsp := g.GetAttr("MaxSP")
	if maxsp != nil {
		g.SetAttr("SP", maxsp)
	}

	rec := g.Record.Record("clone_player_pass_history_record")
	if rec == nil {
		g.Log.Println("not found clone_player_entry_info_record")
		return
	}

	if rec.RowCount() <= 0 {
		g.Log.Println("clone_player_entry_info_record row is zero")
		return
	}
	scene := rec.Row(rand.Intn(rec.RowCount()))
	msg := protocol.NewVarMsg(3)
	msg.AddInt32(protocol.CLINET_CUSTOMMSG_CLONE_SCENE)
	msg.AddInt32(1) //1 进入副本
	msg.AddInt64(scene.Value(0).Int64())
	g.SendCustom(msg)

	g.Log.Println("进入副本", scene.Value(0).Int64())
}

func (g *GameClient) QuitCloneScene() {
	msg := protocol.NewVarMsg(3)
	msg.AddInt32(protocol.CLINET_CUSTOMMSG_CLONE_SCENE)
	msg.AddInt32(6) //6 退出副本
	g.SendCustom(msg)
}

func (g *GameClient) Relive() {
	msg := protocol.NewVarMsg(5)
	msg.AddInt32(protocol.CLIENT_CUSTOMMSG_RELIVE)
	msg.AddInt32(0)   //复活类型
	msg.AddInt32(0)   //是否使用道具
	msg.AddString("") //复活点
	msg.AddInt32(1)   //默认复活点
	g.SendCustom(msg)
}

func (g *GameClient) OnRecordChange(rec string, op int, row, col int) {
	switch rec {
	case "BeRequestRecord":
		g.Request(op, row, col)
	}
}

func (g *GameClient) Request(op int, row, col int) {
	switch op {
	case E_DATAGRID_ADD_ROW:
		rec := g.Record.Record("BeRequestRecord")
		if rec == nil {
			return
		}

		r := rec.Row(row)
		if r == nil {
			return
		}

		src := r.Value(1)
		typ := r.Value(6)
		msg := protocol.NewVarMsg(4)
		msg.AddInt32(protocol.CLIENT_CUSTOMMSG_ANSWER)
		msg.AddInt32(typ.Int())
		msg.AddString(src.String())
		msg.AddInt32(1)
		g.SendCustom(msg)
	}
}

func (g *GameClient) JoinGuild(guild string) {
	if guild == "" {
		return
	}

	guildid, err := strconv.Atoi(guild)
	if err != nil {
		return
	}

	msg := protocol.NewVarMsg(3)
	msg.AddInt32(protocol.CLIENT_CUSTOMMSG_GUILD)
	msg.AddInt32(3)
	msg.AddInt64(int64(guildid))
	g.SendCustom(msg)
}

func (g *GameClient) QuitGuild() {
	msg := protocol.NewVarMsg(2)
	msg.AddInt32(protocol.CLIENT_CUSTOMMSG_GUILD)
	msg.AddInt32(14)
	g.SendCustom(msg)
}

func (g *GameClient) PowerUp() {
	g.SetObj()
	g.SendGM("set MaxHp 999999")
	g.SendGM("set Atk 4999999")
	g.SendGM("set Hp 999999")
}

func (g *GameClient) PowerDown() {
	g.SetObj()
	g.SendGM("set MaxHp 9999")
	g.SendGM("set Atk 9999")
	g.SendGM("set Hp 9999")
}

func (g *GameClient) Enter15V15() {
	msg := protocol.NewVarMsg(2)
	msg.AddInt32(protocol.CLIENT_CUSTOMMSG_GROUP_PK_MSG)
	msg.AddInt32(0)
	g.SendCustom(msg)
}

func (g *GameClient) Quit15V15() {
	msg := protocol.NewVarMsg(2)
	msg.AddInt32(protocol.CLIENT_CUSTOMMSG_GROUP_PK_MSG)
	msg.AddInt32(2)
	g.SendCustom(msg)
}

func (g *GameClient) EnterCloneSceneDetails(sceneid string) {
	g.SetObj()
	maxsp := g.GetAttr("MaxSP")
	if maxsp != nil {
		g.SetAttr("SP", maxsp)
	}

	rec := g.Record.Record("clone_player_pass_history_record")
	if rec == nil {
		g.Log.Println("not found clone_player_entry_info_record")
		return
	}

	if rec.RowCount() <= 0 {
		g.Log.Println("clone_player_entry_info_record row is zero")
		return
	}
	scene, _ := strconv.ParseInt(sceneid, 10, 64)

	msg := protocol.NewVarMsg(3)
	msg.AddInt32(protocol.CLINET_CUSTOMMSG_CLONE_SCENE)
	msg.AddInt32(1) //1 进入副本
	msg.AddInt64(scene)
	g.SendCustom(msg)

	g.Log.Println("进入副本", scene)
}

func (g *GameClient) GetCloneSceneIdList() (list []int64, error string) {
	g.SetObj()
	maxsp := g.GetAttr("MaxSP")
	if maxsp != nil {
		g.SetAttr("SP", maxsp)
	}

	rec := g.Record.Record("clone_player_pass_history_record")
	if rec == nil {
		g.Log.Println("not found clone_player_entry_info_record")
		error = "not found clone_player_entry_info_record"
		return
	}

	for _, row := range rec.Rows {
		list = append(list, row.Value(0).Int64())
	}
	error = ""
	return
}
func (g *GameClient) PK(rguid int64) {
	msg := protocol.NewVarMsg(3)
	msg.AddInt32(protocol.CLIENT_CUSTOMMSG_PK_MSG)
	msg.AddInt32(1)
	msg.AddInt64(rguid)
	g.SendCustom(msg)
	fmt.Println(rguid, "===PK===", g.GetAttr("RGuid"))
}

func (g *GameClient) ResponsePK(rguid int64) {
	msg := protocol.NewVarMsg(4)
	msg.AddInt32(protocol.CLIENT_CUSTOMMSG_PK_MSG)
	msg.AddInt32(2)
	msg.AddInt32(1)
	msg.AddInt64(rguid)
	g.SendCustom(msg)
}

func (g *GameClient) EnterMultiScene(sceneid string) {
	scene, _ := strconv.ParseInt(sceneid, 10, 64)
	msg := protocol.NewVarMsg(3)
	msg.AddInt32(protocol.CLIENT_CUSTOMMSG_MULTISCENE)
	msg.AddInt32(1)
	msg.AddInt32(int32(scene))
	g.SendCustom(msg)
}

func (g *GameClient) ResponseMultiScene() {
	msg := protocol.NewVarMsg(3)
	msg.AddInt32(protocol.CLIENT_CUSTOMMSG_MULTISCENE)
	msg.AddInt32(5)
	msg.AddInt32(1)
	g.SendCustom(msg)
}

func (g *GameClient) CreateTeam() {
	msg := protocol.NewVarMsg(3)
	msg.AddInt32(protocol.CLIENT_CUSTOMMSG_TEAM)
	msg.AddInt32(1)
	msg.AddInt64(0)
	g.SendCustom(msg)
}

func (g *GameClient) ResponseTeamCreate() {
	msg := protocol.NewVarMsg(2)
	msg.AddInt32(protocol.CLIENT_CUSTOMMSG_TEAM)
	msg.AddInt32(7)
	g.SendCustom(msg)
}

func (g *GameClient) AutoMatchTeam() {
	msg := protocol.NewVarMsg(3)
	msg.AddInt32(protocol.CLIENT_CUSTOMMSG_TEAM)
	msg.AddInt32(9)
	msg.AddInt64(0)
	g.SendCustom(msg)
}

func (g *GameClient) SceneMove(custom string) {
	args := strings.Split(custom, " ")
	if len(args) < 1 {
		return
	}
	sceneids := strings.Split(args[0], ",")
	duration, _ := strconv.Atoi(args[1])
	sceneid := sceneids[rand.Intn(len(sceneids))]
	g.AddTimer("scene_move", func(args interface{}) { g.SceneMoveDetails(sceneid) }, nil, time.Minute*time.Duration(duration), -1)

}
func (g *GameClient) SceneMoveDetails(custom string) {
	g.SwitchScene(custom)
	g.BeginMove()
}

func (g *GameClient) StopSceneMove() {
	g.RemoveTimer("scene_move")
	g.StopMove()
}
func (g *GameClient) SubmitQuest() {

	rec := g.Record.Record("QuestRecord")
	if rec == nil {
		g.Log.Println("not found QuestRecord")
		return
	}
	if rec.RowCount() <= 0 {
		g.Log.Println("QuestRecord row is zero")
		return

	}
	row := rec.Row(0)
	g.SetObj()
	g.SendGM(fmt.Sprintf("submit_quest %d", row.Value(0).Int64()))

}
func (g *GameClient) AcceptQuest() {
	rec := g.Record.Record("QuestAcceptRecord")
	if rec == nil {
		g.Log.Println("not found QuestAcceptRecord")
		return
	}
	if rec.RowCount() <= 0 {
		g.Log.Println("QuestAcceptRecord row is zero")
		return
	}
	row := rec.Row(0)
	g.SetObj()
	g.SendGM(fmt.Sprintf("accept_quest %d", row.Value(0).Int64()))
}
func (g *GameClient) NinjaArena() {
	msg := protocol.NewVarMsg(2)
	msg.AddInt32(protocol.CLIENT_CUSTOMMSG_NINJA_ARENA_MSG)
	msg.AddInt32(0)
	g.SendCustom(msg)
}
func (g *GameClient) ResponseNinjaArena(uid int64) {
	msg := protocol.NewVarMsg(3)
	msg.AddInt32(protocol.CLIENT_CUSTOMMSG_NINJA_ARENA_MSG)
	msg.AddInt32(3)
	msg.AddInt64(uid)
	g.SendCustom(msg)
}

func (g *GameClient) StopNinjaArena() {
	msg := protocol.NewVarMsg(2)
	msg.AddInt32(protocol.CLIENT_CUSTOMMSG_NINJA_ARENA_MSG)
	msg.AddInt32(7)
	g.SendCustom(msg)
}

func (g *GameClient) SendCustomMsg(custom string) {
	defer func() {
		if err := recover(); err != nil {
			g.Log.Println(err)
		}
	}()
	args := strings.Split(custom, " ")
	if len(args) < 1 {
		return
	}
	count := (len(args) - 1) / 2
	msg := protocol.NewVarMsg(count + 1)
	msgid, err := strconv.Atoi(args[0])
	if err != nil {
		panic(err)
	}
	msg.AddInt32(int32(msgid))
	for i := 1; i < len(args); i += 2 {
		switch args[i] {
		case "i":
			value, err := strconv.ParseInt(args[i+1], 10, 32)
			if err != nil {
				panic(err)
			}
			msg.AddInt32(int32(value))
		case "i64":
			value, err := strconv.ParseInt(args[i+1], 10, 64)
			if err != nil {
				panic(err)
			}
			msg.AddInt64(value)
		case "f":
			value, err := strconv.ParseFloat(args[i+1], 32)
			if err != nil {
				panic(err)
			}
			msg.AddFloat(float32(value))
		case "d":
			value, err := strconv.ParseFloat(args[i+1], 64)
			if err != nil {
				panic(err)
			}
			msg.AddDouble(value)
		case "s":
			msg.AddString(args[i+1])
		case "o":
			value, err := strconv.ParseInt(args[i+1], 10, 64)
			if err != nil {
				panic(err)
			}
			msg.AddObject(uint64(value))
		}
	}
	g.SendCustom(msg)
}

func (g *GameClient) GroupPK() {
	g.SendGM("setobj")
	g.SendGM("reset40")
	msg := protocol.NewVarMsg(2)
	msg.AddInt32(protocol.CLIENT_CUSTOMMSG_GROUP_PK_MSG)
	msg.AddInt32(0) // 报名
	g.SendCustom(msg)
}

func (g *GameClient) TeamPK() {
	g.SendGM("setobj")
	g.SendGM("reset5")
	msg := protocol.NewVarMsg(2)
	msg.AddInt32(protocol.CLIENT_CUSTOMMSG_TEAM_PK_MSG)
	msg.AddInt32(0)
	g.SendCustom(msg)
}
