package server

import (
	"fmt"
	"io"
	"login/protocol"
	"net"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type TextProtocol struct {
	ctx *Context
}

func HexToChar(s []byte, start, count int) byte {
	str := string(s[start : start+count])
	val, _ := strconv.ParseInt(str, 16, 32)
	return byte(val)
}

func CharToHex(ch byte, bytes int) []byte {
	f := fmt.Sprintf("%%%dX", bytes)
	str := fmt.Sprintf(f, ch)
	return []byte(str)
}

func (p *TextProtocol) IOLoop(conn net.Conn) error {
	var zeroTime time.Time

	clientId := atomic.AddInt64(&p.ctx.server.clientIDSequence, 1)
	client := newClient(clientId, conn, p.ctx)
	if !p.ctx.server.AddClient(client) {
		return fmt.Errorf("add client(%d) error", clientId)
	}

	go p.messagePump(client)

	for {
		if client.HeartbeatInterval > 0 {
			client.SetReadDeadline(time.Now().Add(client.HeartbeatInterval * 2))
		} else {
			client.SetReadDeadline(zeroTime)
		}

		buf, err := client.Reader.ReadSlice(byte(0x0A))
		if err != nil {
			if !strings.Contains(err.Error(), "An existing connection was forcibly closed by the remote host") {
				p.ctx.server.log.Print(err.Error())
			}
			break
		}

		if len(buf) <= 2 {
			p.ctx.server.log.Println("msg size error")
			continue
		}

		data := buf[:len(buf)-2]
		msg := protocol.NewVarMsg(16)
		if p.DecodeMsg(data, len(data), msg) {
			serialid := atomic.AddInt64(&p.ctx.server.msgIdSequence, 1)
			msg.ConnId = client.Id
			msg.Serial = serialid
			p.Exec(client, msg)
		}
	}

	client.Quit()
	RemoveServer(client, nil)
	p.ctx.server.RemoveClient(client.Id)
	p.ctx.server.RemoveAllUserByConnid(client.Id)
	return nil
}

func (p *TextProtocol) decodeData(buf []byte, start, end int, msg *protocol.VarMessage) bool {
	if start == end {
		return true
	}

	first := buf[start]
	len := end - start
	if first == '#' { //widestr
		p.ctx.server.log.Println("unsolved widestr")
	} else if first == '*' { //binary
		p.ctx.server.log.Println("unsolved binary")
	} else if first == '$' { //string
		s := make([]byte, 0, len)
		pos := start + 1
		for pos < end {
			if buf[pos] == '\\' {
				pos++
				if pos >= end {
					return false
				}
				if buf[pos] == '\\' {
					pos++
					s = append(s, '\\')
					continue
				}

				if buf[pos] == 'x' {
					pos++
					if pos+1 >= end {
						return false
					}

					c := HexToChar(buf, pos, 2)
					s = append(s, c)
					pos = pos + 2
					continue
				}

				return false
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
				return false
			}
			msg.AddDouble(f)
		} else {
			i, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return false
			}
			msg.AddInt64(i)
		}
	}

	return true
}

func (p *TextProtocol) DecodeMsg(buf []byte, size int, msg *protocol.VarMessage) bool {
	msg.Clear()
	beginpos := 0
	for k, v := range buf {
		if v == ' ' {
			if !p.decodeData(buf, beginpos, k, msg) {
				p.ctx.server.log.Println("decode msg error1")
				return false
			}

			beginpos = k + 1
		}
	}

	if beginpos < size {
		if !p.decodeData(buf, beginpos, size, msg) {
			p.ctx.server.log.Println("decode msg error2")
			return false
		}
	}

	return true
}

func (p *TextProtocol) EncodeMsg(w io.Writer, msg *protocol.VarMessage) (int, error) {
	sum := 0
	for i := 0; i < msg.Size; i++ {
		switch msg.Type(i) {
		case protocol.VTYPE_INT, protocol.VTYPE_INT64:
			v := fmt.Sprintf("%d", msg.RawValue(i))
			n, err := w.Write([]byte(v))
			if err != nil {
				return sum, err
			}
			sum += n
		case protocol.VTYPE_FLOAT, protocol.VTYPE_DOUBLE:
			v := fmt.Sprintf("%f", msg.RawValue(i))
			n, err := w.Write([]byte(v))
			if err != nil {
				return sum, err
			}
			sum += n
		case protocol.VTYPE_STRING:
			v := []byte(msg.StringVal(i))
			n, err := w.Write([]byte{'$'})
			if err != nil {
				return sum, err
			}
			sum += n
			for _, ch := range v {
				if ch > 0x20 && ch < 0x7F && ch != 0x25 {
					n, err := w.Write([]byte{ch})
					if err != nil {
						return sum, err
					}
					sum += n
					if ch == '\\' {
						n, err := w.Write([]byte{ch})
						if err != nil {
							return sum, err
						}
						sum += n
					}
					continue
				}

				n, err := w.Write([]byte{'\\', 'x'})
				if err != nil {
					return sum, err
				}
				sum += n
				n, err = w.Write(CharToHex(ch, 2))
				if err != nil {
					return sum, err
				}
				sum += n
			}
		default:
			p.ctx.server.log.Println("unsolved type")
		}
		if i < msg.Size-1 {
			n, err := w.Write([]byte{' '}) //分隔符
			if err != nil {
				return sum, err
			}
			sum += n
		}
	}

	n, err := w.Write([]byte{0x0D, 0x0A})
	if err != nil {
		return sum, err
	}
	sum += n
	return sum, nil
}

func (p *TextProtocol) messagePump(client *Client) {
	for {
		select {
		case m := <-client.sendqueue:
			n, err := p.EncodeMsg(client.Writer, m)
			if err != nil {
				p.ctx.server.log.Println("write message error,", err)
			}
			if err := client.Writer.Flush(); err != nil {
				p.ctx.server.log.Println("flush message error")
			}
			p.ctx.server.log.Printf("send message to %d , size:%d", client.Id, n)
		case <-client.exitchan:
			goto exit
		}
	}

exit:
	p.ctx.server.log.Println("client quit loop")
}

func (p *TextProtocol) Exec(client *Client, msg *protocol.VarMessage) {
	cmd := msg.StringVal(1)
	switch cmd {
	case "register":
		RegisterServer(client, msg)
	case "unregister":
		UnregisterServer(client, msg)
	case "keep":
		//保持
	default:
		p.ctx.server.access.PushMessage(msg)
	}
}
