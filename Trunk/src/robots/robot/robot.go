package robot

import (
	"bufio"
	"bytes"
	"container/list"
	"fmt"
	"log"
	"net"
	"os"
	"robots/protocol"
	"robots/utils"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	ROBOT_STATE_NONE = iota
	ROBOT_STATE_ERROR
	ROBOT_STATE_FAILED
	ROBOT_STATE_DISCONNECTED
	ROBOT_STATE_CONNECTED
	ROBOT_STATE_CREATING
	ROBOT_STATE_CHOOSING
	ROBOT_STATE_READY
)

type GameRobot interface {
	OnConnected()
	OnDisconnected()
	OnFailed()
	OnReceive(msgid uint8, ar *utils.LoadArchive)
	OnStateChange(state, old int)
	OnDestroy()
	OnExec()
	OnReady()
}

type Robot struct {
	sync.RWMutex
	Addr     string
	Port     int
	Index    int
	ServerId string
	Account  string
	Password string
	Name     string
	state    int
	Err      error
	quit     bool
	running  bool
	client   *Client
	gb       GameRobot
	timelist *list.List
	Log      *log.Logger
	msgQueue chan *protocol.Message
}

type NullLog struct {
}

func (l *NullLog) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (r *Robot) Init() {
	r.timelist = list.New()
	//r.Log = log.New(&NullLog{}, "", log.LstdFlags)
	r.Log = log.New(os.Stdout, "", log.LstdFlags)
	r.Log.SetPrefix(fmt.Sprintf("[robot%d]", r.Index))
	r.msgQueue = make(chan *protocol.Message, 32)
}

func (r *Robot) Running() bool {
	return r.running
}

func (r *Robot) AddTimer(name string, f TimeCB, args interface{}, delay time.Duration, count int) {
	for ele := r.timelist.Front(); ele != nil; ele = ele.Next() {
		if ele.Value.(*Timer).Name == name {
			return
		}
	}
	t := &Timer{}
	t.Name = name
	t.delay = delay
	t.count = count
	t.callback = f
	t.args = args
	t.time = time.Now().Add(delay)
	r.timelist.PushBack(t)
}

func (r *Robot) RemoveTimer(name string) {
	for ele := r.timelist.Front(); ele != nil; {
		next := ele.Next()
		if ele.Value != nil && ele.Value.(*Timer).Name == name {
			r.timelist.Remove(ele)
		}
		ele = next
	}
}

func (r *Robot) RemoveAllTimer() {
	for ele := r.timelist.Front(); ele != nil; {
		next := ele.Next()
		r.timelist.Remove(ele)
		ele = next
	}
}

func (r *Robot) execTimer(now time.Time) {
	for ele := r.timelist.Front(); ele != nil; {
		if r.quit {
			return
		}
		next := ele.Next()
		if ele.Value.(*Timer).Exec(now) {
			r.timelist.Remove(ele)
		}
		ele = next
	}
}

func (r *Robot) SetGameRobot(gb GameRobot) {
	r.gb = gb
}

func (r *Robot) ChangeState(state int) {
	old := r.state
	if old != state {
		r.state = state
		if r.gb != nil {
			r.gb.OnStateChange(state, old)
		}
	}
}

func (r *Robot) State() int {
	return r.state
}

func (r *Robot) Close() {
	if r.state >= ROBOT_STATE_CONNECTED && r.client != nil {
		r.client.Shutdown()
		r.client = nil
		r.ChangeState(ROBOT_STATE_DISCONNECTED)
	}
}

func (r *Robot) Connect(addr string, port int, serverid string) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		r.Err = err
		r.ChangeState(ROBOT_STATE_FAILED)
		if r.gb != nil {
			r.gb.OnFailed()
			return
		}
	}

	_, err = conn.Write([]byte("svrlist\r\n"))
	if err != nil {
		r.Err = err
		r.ChangeState(ROBOT_STATE_FAILED)
		if r.gb != nil {
			r.gb.OnFailed()
		}
	}
	go r.parseGateInfo(conn, addr, serverid)
}

func (r *Robot) parseGateInfo(conn net.Conn, addr string, serverid string) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	buf, err := reader.ReadSlice(byte(0x0A))
	if err != nil {
		r.Err = err
		r.ChangeState(ROBOT_STATE_FAILED)
		if r.gb != nil {
			r.gb.OnFailed()
			return
		}
	}

	msg := protocol.NewVarMsg(16)
	err = DecodeMsg(buf[:len(buf)-2], len(buf)-2, msg)
	if err != nil {
		r.Err = err
		r.ChangeState(ROBOT_STATE_FAILED)
		if r.gb != nil {
			r.gb.OnFailed()
			return
		}
	}

	typ := msg.StringVal(0)
	if typ == "svrlist" {
		//addr := msg.StringVal(1)
		port := msg.Int64Val(2)
		r.Connect2(addr, int(port), serverid)
	}
}

func HexToChar(s []byte, start, count int) byte {
	str := string(s[start : start+count])
	val, _ := strconv.ParseInt(str, 16, 32)
	return byte(val)
}

func DecodeMsg(buf []byte, size int, msg *protocol.VarMessage) error {
	msg.Clear()
	beginpos := 0
	for k, v := range buf {
		if v == ' ' {
			if err := DecodeData(buf, beginpos, k, msg); err != nil {
				return err
			}
			beginpos = k + 1
		}
	}

	if beginpos < size {
		if err := DecodeData(buf, beginpos, size, msg); err != nil {
			return err
		}
	}

	return nil
}

func DecodeData(buf []byte, start, end int, msg *protocol.VarMessage) error {
	if start == end {
		return fmt.Errorf("msg is nil")
	}

	first := buf[start]
	len := end - start
	if first == '#' { //widestr
		return fmt.Errorf("unsolved widestr")
	} else if first == '*' { //binary
		return fmt.Errorf("unsolved binary")
	} else if first == '$' { //string
		s := make([]byte, 0, len)
		pos := start + 1
		for pos < end {
			if buf[pos] == '\\' {
				pos++
				if pos >= end {
					return fmt.Errorf("msg error")
				}
				if buf[pos] == '\\' {
					pos++
					s = append(s, '\\')
					continue
				}

				if buf[pos] == 'x' {
					pos++
					if pos+1 >= end {
						return fmt.Errorf("msg  error")
					}

					c := HexToChar(buf, pos, 2)
					s = append(s, c)
					pos = pos + 2
					continue
				}

				return fmt.Errorf("msg error")
			}

			s = append(s, buf[pos])
			pos++
		}

		msg.AddString(string(s))
	} else { //number
		val := string(buf[start:end])
		if strings.IndexByte(val, '.') != -1 {
			f, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return fmt.Errorf("msg parse float error")
			}
			msg.AddDouble(f)
		} else {
			i, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return fmt.Errorf("msg parse int error")
			}
			msg.AddInt64(i)
		}
	}

	return nil
}

func (r *Robot) Connect2(addr string, port int, serverid string) {
	if r.state >= ROBOT_STATE_CONNECTED && r.client != nil {
		r.client.Shutdown()
		r.client = nil
		r.ChangeState(ROBOT_STATE_DISCONNECTED)
	}

	r.Err = nil
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		r.Err = err
		r.ChangeState(ROBOT_STATE_FAILED)
		if r.gb != nil {
			r.gb.OnFailed()
		}
		return
	}
	r.Addr = addr
	r.Port = port
	r.ServerId = serverid
	client := newClient(conn)
	r.client = client
	r.ChangeState(ROBOT_STATE_CONNECTED)
	if r.gb != nil {
		r.gb.OnConnected()
	}
	go r.ioLoop(client)
	go r.sendLoop(client)
}

func (r *Robot) Run() {
	var now time.Time
	r.running = true
	for !r.quit {
		now = time.Now()
		r.processMsg()
		r.execTimer(now)
		if r.gb != nil && !r.quit {
			r.gb.OnExec()
		}
		time.Sleep(time.Millisecond)
	}

	r.shutdown()
	r.running = false
}

func (r *Robot) Destroy() {
	r.quit = true
}

func (r *Robot) shutdown() {
	if r.client != nil {
		r.client.Shutdown()
		r.client = nil
	}

	if r.gb != nil {
		r.gb.OnDestroy()
		r.gb = nil
	}

	close(r.msgQueue)

}

func (r *Robot) ioLoop(client *Client) {
	buf := bytes.NewBuffer(nil)
	var prv byte
	count := 0
	for !r.quit && !client.quit {
		ch, err := client.Reader.ReadByte()
		if err != nil {
			client.Shutdown()
			r.Err = err
			r.ChangeState(ROBOT_STATE_DISCONNECTED)
			break
		}

		if prv == 0xEE && ch == 0xEE {
			count--
			buf.Truncate(count)
			r.OnReceive(buf.Bytes())
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

func (r *Robot) sendLoop(client *Client) {
loop:
	for !r.quit && !client.quit {
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
		case <-client.exitchan:
			break loop
		}
	}
}

func (r *Robot) OnReceive(data []byte) {
	msg := protocol.NewMessage(len(data))
	msg.Body = append(msg.Body, data...)
	r.msgQueue <- msg
}

func (r *Robot) processMsg() {
	for {
		select {
		case m := <-r.msgQueue:
			if r.quit {
				return
			}
			ar := utils.NewLoadArchiver(m.Body)
			msgtype, err := ar.ReadUInt8()
			if err != nil {
				panic(err)
			}

			if r.gb != nil {
				r.gb.OnReceive(msgtype, ar)
			}
			m.Free()
		default:
			return
		}
	}

}

func (r *Robot) Login(account, password string, login_string string, login_type int32, srv_id string) {
	msg := protocol.NewMessage(512)
	ar := utils.NewStoreArchiver(msg.Body)
	ar.Write(uint8(protocol.CTOS_LOGIN))
	// version
	ar.Write(int32(0x31303030))
	// account
	ar.Write(account)
	// password
	ar.Write(password)
	// loginstring
	ar.Write(login_string)
	// login type
	ar.Write(login_type)
	// dev no
	ar.Write("AB-CD-EF-GH-IJ-KL")
	// pt flg
	ar.Write("101")
	// server id
	ar.Write(srv_id)
	msg.Body = msg.Body[:ar.Len()]
	r.SendMessage(msg)
}

func (r *Robot) CreateRole(role_index int, args *protocol.VarMessage) {
	msg := protocol.NewMessage(512)
	ar := utils.NewStoreArchiver(msg.Body)
	ar.Write(uint8(protocol.CTOS_CREATE_ROLE))
	ar.Write(int32(role_index))
	ar.Write(uint8(0))
	ar.Write(r.Name)
	//ar.Write("€Tv🤘")
	protocol.PubArgs(ar, args)
	msg.Body = msg.Body[:ar.Len()]
	r.SendMessage(msg)
	r.Log.Println(r.Account, "create role", r.Name)
	r.ChangeState(ROBOT_STATE_CREATING)
}

func (r *Robot) ChooseRole(role_name string) {
	msg := protocol.NewMessage(256)
	ar := utils.NewStoreArchiver(msg.Body)
	ar.Write(uint8(protocol.CTOS_CHOOSE_ROLE))
	ar.Write(role_name)
	ar.Write("")
	msg.Body = msg.Body[:ar.Len()]
	r.SendMessage(msg)
	r.Log.Println(r.Account, "choose role", role_name)
	r.ChangeState(ROBOT_STATE_CHOOSING)
}

func (r *Robot) Ready() {
	msg := protocol.NewMessage(16)
	ar := utils.NewStoreArchiver(msg.Body)
	ar.Write(uint8(protocol.CTOS_READY))
	msg.Body = msg.Body[:ar.Len()]
	r.SendMessage(msg)
	if r.gb != nil {
		r.gb.OnReady()
	}
	r.Log.Println(r.Account, r.Name, "ready")
	r.ChangeState(ROBOT_STATE_READY)
}

func (r *Robot) ReqMove(mode uint8, pos []float32, info string) {
	varlist := protocol.NewVarMsg(7)
	varlist.AddInt32(protocol.CLIENT_CUSTOMMSG_REQUEST_MOVE)
	varlist.AddInt32(int32(mode))
	for _, v := range pos {
		varlist.AddFloat(v)
	}
	r.SendCustom(varlist)
}

func (r *Robot) SendCustom(args *protocol.VarMessage) {
	msg := protocol.NewMessage(512)
	ar := utils.NewStoreArchiver(msg.Body)
	ar.Write(uint8(protocol.CTOS_CUSTOM))
	protocol.PubArgs(ar, args)
	msg.Body = msg.Body[:ar.Len()]
	r.SendMessage(msg)
	//r.Log.Println("send custom")
}

func (r *Robot) SendMessage(msg *protocol.Message) {
	if r.state >= ROBOT_STATE_CONNECTED && r.client != nil {
		r.client.sendqueue <- msg
	}
}
