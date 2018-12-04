package server

import (
	"charge/protocol"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// 文本协议
type TextProtocol struct {
	ctx *Context
}

// 十六进制字符串转换
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

	clientId := atomic.AddInt64(&p.ctx.Server.clientIDSequence, 1)
	client := newClient(clientId, conn, p.ctx)
	if !p.ctx.Server.AddClient(client) {
		return fmt.Errorf("add client(%d) error", clientId)
	}

	// 启动消息发送队列
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
				p.ctx.Server.log.Print(err.Error())
			}
			break
		}

		if len(buf) <= 2 {
			p.ctx.Server.log.Println("msg size error")
			continue
		}

		data := buf[:len(buf)-2]
		msg := protocol.NewVarMsg(16)
		if p.DecodeMsg(data, len(data), msg) {
			serialid := atomic.AddInt64(&p.ctx.Server.msgIdSequence, 1)
			msg.ConnId = client.Id
			msg.Serial = serialid
			p.Exec(client, msg)
		}
	}

	p.ctx.Server.log.Println("lost client", client.RemoteAddr())
	client.Quit()
	RemoveServer(client, nil)
	p.ctx.Server.RemoveClient(client.Id)
	return nil
}

// 消息解码
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

// 消息解码
func (p *TextProtocol) DecodeMsg(buf []byte, size int, msg *protocol.VarMessage) bool {
	msg.Clear()
	beginpos := 0
	for k, v := range buf {
		if v == ' ' {
			if err := DecodeData(buf, beginpos, k, msg); err != nil {
				p.ctx.Server.log.Println(err.Error())
				return false
			}

			beginpos = k + 1
		}
	}

	if beginpos < size {
		if err := DecodeData(buf, beginpos, size, msg); err != nil {
			p.ctx.Server.log.Println(err.Error())
			return false
		}
	}

	return true
}

// 消息编码
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
			p.ctx.Server.log.Println("unsolved type")
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

// 消息发送线程
func (p *TextProtocol) messagePump(client *Client) {
	for {
		select {
		case m := <-client.sendqueue:
			n, err := p.EncodeMsg(client.Writer, m)
			if err != nil {
				p.ctx.Server.log.Println("write message error,", err)
			}
			if err := client.Writer.Flush(); err != nil {
				p.ctx.Server.log.Println("flush message error")
			}
			p.ctx.Server.log.Printf("send message to %d , size:%d", client.Id, n)
		case <-client.exitchan:
			goto exit
		}
	}

exit:
	p.ctx.Server.log.Println("client quit loop")
}

// 消息处理函数
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
		p.ctx.Server.access.PushMessage(msg)
	}
}
