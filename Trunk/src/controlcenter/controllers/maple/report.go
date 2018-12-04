package maple

import (
	"controlcenter/controllers/utils"
	"controlcenter/setting"
	"fmt"
	"net"
	"os"

	mgo "gopkg.in/mgo.v2"

	"bytes"
	"encoding/binary"

	"github.com/astaxie/beego"
	simplejson "github.com/bitly/go-simplejson"
)

const (
	MAX_UDP_LEN = 0x1000
	MAPLEDB     = "maple"
)

var (
	globalSession *mgo.Session
)

type Handler struct {
	data    []byte
	session *mgo.Session
	DB      *mgo.Database
}

func (m *Handler) Close() {
	m.session.Close()
}

func GetDBSession() *mgo.Session {
	if globalSession != nil {
		return globalSession.Clone()
	}

	return nil
}

func NewHandler(data []byte) *Handler {
	m := &Handler{}
	m.data = data
	m.session = globalSession.Clone()
	m.DB = m.session.DB(MAPLEDB)
	return m
}

func HandleMessage(conn *net.UDPConn) {
	buff := make([]byte, MAX_UDP_LEN)
	for {
		n, _, err := conn.ReadFromUDP(buff[:])
		if err != nil {
			beego.Error("read udp package error", err)
			continue
		}

		if n == 0 {
			continue
		}

		r := bytes.NewReader(buff[:2])
		var size uint16
		err = binary.Read(r, binary.LittleEndian, &size)
		if err != nil {
			beego.Error("read size error", err)
			continue
		}
		if n-2 != int(size) {
			beego.Error("size not match:", n-2, size)
			continue
		}

		msg := utils.NewMessage(int(size))
		msg.Body = append(msg.Body, buff[2:]...)
		go ParseMessage(msg)
	}

}

func ParseMessage(msg *utils.Message) {
	jsonobj, err := simplejson.NewJson(msg.Body)
	if err != nil {
		beego.Error("parse message failed,", err)
		beego.Error(string(msg.Body))
		msg.Free()
		return
	}
	msg.Free()

	msgtype := jsonobj.Get("msg_type").MustString("")
	data, err := jsonobj.Get("user_data").Encode()
	if err != nil {
		beego.Error(err)
		return
	}

	handler := NewHandler(data)
	switch msgtype {
	case "report_status":
		ReportServerStatus(handler)
	case "user_role_info":
		ReportRoleInfo(handler)
	}
	handler.Close()
}

func StartUdpServer(addr string, port int) *net.UDPConn {
	udpaddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		beego.Error("resolve udp addr error, ", err)
		os.Exit(1)
	}

	beego.Info(udpaddr)
	conn, err := net.ListenUDP("udp4", udpaddr)
	if err != nil {
		beego.Error("listen udp failed, ", err)
		os.Exit(1)
	}
	beego.Info("start status server at:", addr, ":", port)

	globalSession, err = mgo.Dial(setting.StatusDB)
	if err != nil {
		beego.Error("connect to mango db failed, ", err)
		os.Exit(1)
	}
	globalSession.SetMode(mgo.Monotonic, true)
	go HandleMessage(conn)
	return conn
}
