package sample

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"sync"
	"systemmonitor/protocol"
	"systemmonitor/utils"
	"time"
)

type GameInfo struct {
	sync.RWMutex
	addr     string
	port     int
	Index    int
	Err      error
	client   *Client
	msgQueue chan *protocol.Message
	sample   *Sample
}

func NewGameInfo(s *Sample) *GameInfo {
	g := &GameInfo{}
	g.Init()
	g.sample = s
	return g
}

func (g *GameInfo) Init() {
	g.msgQueue = make(chan *protocol.Message, 32)
}

func (g *GameInfo) SetAddrPort(addr string, port int) {
	if g.addr != addr || g.port != port {
		if g.client != nil {
			g.client.Shutdown()
			g.client = nil
		}

	}

	g.addr = addr
	g.port = port
}

func (g *GameInfo) KeepConnection(ctx context.Context) {
	if g.client != nil {
		return
	}

	g.Err = nil
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", g.addr, g.port))
	if err != nil {
		fmt.Println("connect game failed", err)
		g.Err = err
		return
	}
	client := newClient(conn)
	g.client = client

	go g.ioLoop(ctx, client)
	go g.sendLoop(ctx, client)

	fmt.Println("game connected")
}

func (g *GameInfo) Run(ctx context.Context) {
	t := time.Tick(time.Second * 5)          // 每5秒尝试连接一次
	querytime := time.Tick(time.Second * 10) // 每10秒查询一次服务器信息
L:
	for {
		select {
		case <-ctx.Done():
			break L
		case <-t:
			g.KeepConnection(ctx)
		case <-querytime:
			g.QuerySceneInfo()
			g.QueryOnlineInfo()
			g.QueryRegPlayers()
		default:
			g.processMsg()
			time.Sleep(time.Millisecond)
		}
	}

	g.shutdown()
}

func (g *GameInfo) QuerySceneInfo() {
	if g.client != nil {
		varlist := protocol.NewVarMsg(7)
		varlist.AddString("scene_players")
		g.SendMessage(varlist)
	}
}

func (g *GameInfo) QueryOnlineInfo() {
	if g.client != nil {
		varlist := protocol.NewVarMsg(7)
		varlist.AddString("online_players")
		g.SendMessage(varlist)
	}
}

func (g *GameInfo) QueryRegPlayers() {
	if g.client != nil {
		varlist := protocol.NewVarMsg(7)
		varlist.AddString("reg_players")
		g.SendMessage(varlist)
	}
}

func (g *GameInfo) shutdown() {
	if g.client != nil {
		g.client.Shutdown()
		g.client = nil
	}

	close(g.msgQueue)
}

func (g *GameInfo) ioLoop(ctx context.Context, client *Client) {
	buf := bytes.NewBuffer(nil)
	var prv byte
	count := 0
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		ch, err := client.Reader.ReadByte()
		if err != nil {
			client.Shutdown()
			g.Err = err
			break
		}

		if prv == 0xEE && ch == 0xEE {
			count--
			buf.Truncate(count)
			g.OnReceive(buf.Bytes())
			buf.Reset()
			count = 0
			continue
		} else if ch == 0 && prv == 0xEE {

		} else {
			buf.WriteByte(ch)
			count++
		}

		prv = ch
	}
}

func (g *GameInfo) sendLoop(ctx context.Context, client *Client) {
loop:
	for {
		select {
		case m := <-client.sendqueue:
			for _, b := range m.Body {
				client.Writer.WriteByte(b)
				if b == 0xEE {
					client.Writer.WriteByte(0)
				}
			}
			client.Writer.WriteByte(0xEE)
			client.Writer.WriteByte(0xEE)
			client.Writer.Flush()
			m.Free()
		case <-ctx.Done():
			break loop
		}
	}
}

func (g *GameInfo) OnReceive(data []byte) {
	if len(data) > 0 {
		msg := protocol.NewMessage(len(data))
		msg.Body = append(msg.Body, data...)
		g.msgQueue <- msg
	}
}

func (g *GameInfo) processMsg() {
	for {
		select {
		case m, ok := <-g.msgQueue:
			if ok {
				ar := utils.NewLoadArchiver(m.Body)
				msg := protocol.ParseArgs(ar)
				g.exec(msg)
				m.Free()
			}
		default:
			return
		}
	}
}

// 发送消息
func (g *GameInfo) SendMessage(args *protocol.VarMessage) {
	msg := protocol.NewMessage(512)
	ar := utils.NewStoreArchiver(msg.Body)
	protocol.PubArgs(ar, args)
	msg.Body = msg.Body[:ar.Len()]
	if g.client != nil {
		g.client.sendqueue <- msg
	}
}

func (g *GameInfo) exec(msg *protocol.VarMessage) {
	if msg.Size == 0 {
		return
	}

	if msg.Type(0) != protocol.E_VTYPE_STRING {
		return
	}

	msgtype := msg.StringVal(0)
	switch msgtype {
	case "scene_players":
		g.ScenePlayers(msg)
	case "online_players":
		g.OnlinePlayers(msg)
	case "reg_players":
		g.RegPlayers(msg)
	default:
		fmt.Println("unknown msg")
	}
}

func (g *GameInfo) ScenePlayers(msg *protocol.VarMessage) {

	k := 1
	count := msg.Int32Val(k)
	sceneinfo := make([]*SceneInfo, 0, count)
	k++
	for i := 0; i < int(count); i++ {
		scene := &SceneInfo{}
		scene.SceneId = msg.Int32Val(k)
		k++
		scene.Players = msg.Int32Val(k)
		k++
		scene.SceneName = msg.StringVal(k)
		k++
		sceneinfo = append(sceneinfo, scene)
	}

	g.sample.StoreSceneInfo(sceneinfo)
}

func (g *GameInfo) OnlinePlayers(msg *protocol.VarMessage) {
	count := msg.Int32Val(1)
	g.sample.StoreOnlinePlayer(count)
}

func (g *GameInfo) RegPlayers(msg *protocol.VarMessage) {
	count := msg.Int32Val(1)
	g.sample.StoreMaxPlayer(count)
}
