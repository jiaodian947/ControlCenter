package server

import (
	"encoding/binary"
	"io"
	"speedy/protocol"
	"strings"
	"time"
)

const (
	MAX_MSG_LEN = 256 * 1024 * 1024
)

type BinaryProtocol struct {
	client *Client
}

func (p *BinaryProtocol) IOLoop() error {
	var zeroTime time.Time

	var lens int32
	data := make([]byte, MAX_MSG_LEN)

	go p.messagePump()

	for {
		p.client.SetReadDeadline(zeroTime)

		err := binary.Read(p.client.Reader, binary.LittleEndian, &lens)
		if err != nil {
			if !strings.Contains(err.Error(), "An existing connection was forcibly closed by the remote host") {
				ServerApp.log.Print(err.Error())
			}
			break
		}

		if lens > MAX_MSG_LEN {
			ServerApp.log.Println("msg size error")
			continue
		}

		data = data[:lens]
		_, err = io.ReadFull(p.client.Reader, data)
		if err != nil {
			if !strings.Contains(err.Error(), "An existing connection was forcibly closed by the remote host") {
				ServerApp.log.Print(err.Error())
			}
			break
		}

		msg := protocol.NewMessage(len(data))
		msg.Body = append(msg.Body, data...)
		p.Exec(msg)
		msg.Free()
	}

	return nil
}

func (p *BinaryProtocol) Exec(msg *protocol.Message) {
	ar := protocol.NewLoadArchiver(msg.Body)
	var seq uint64
	err := ar.Read(&seq)
	if err != nil {
		panic(err)
	}

	if seq != 0 {
		p.client.Response(seq, msg)
	}
}

func (p *BinaryProtocol) EncodeMsg(w io.Writer, msg *protocol.Message) (int, error) {
	if err := binary.Write(w, binary.LittleEndian, uint32(len(msg.Body))+8); err != nil {
		return 0, err
	}
	n, err := w.Write(msg.Header)
	if n != 8 || err != nil {
		return n, err
	}
	return w.Write(msg.Body)
}

func (p *BinaryProtocol) messagePump() {
	timeout := time.NewTimer(p.client.Timeout)
loop:
	for {
		if !timeout.Stop() {
			<-timeout.C
		}
		timeout.Reset(p.client.Timeout)
		select {
		case <-timeout.C:
			timeout.Stop()
			p.client.Quit()
			ServerApp.log.Println("idle timeout close connection")
			break loop
		case m := <-p.client.sendqueue:
			n, err := p.EncodeMsg(p.client.Writer, m)
			m.Free()
			if err != nil {
				ServerApp.log.Println("write message error,", err)
				break loop
			}
			if err := p.client.Writer.Flush(); err != nil {
				ServerApp.log.Println("flush message error")
				break loop
			}
			ServerApp.log.Printf("send message to %d , size:%d", p.client.Id, n)
		case <-p.client.exitchan:
			break loop
		}
	}

	ServerApp.log.Println("client quit loop")
}
