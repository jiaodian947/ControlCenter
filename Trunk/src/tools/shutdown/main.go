package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

func main() {

	conn, err := net.Dial("tcp", "127.0.0.1:15001")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	var buf [4096]byte
	n, err := conn.Read(buf[0:])
	if err != nil {
		panic(err)
	}
	fmt.Println(string(buf[:n]))
	w := bufio.NewWriter(conn)
	time.Sleep(time.Second)
	w.Write([]byte("srv_exit\n"))
	w.Flush()
	fmt.Println("srv_exit")
	n, err = conn.Read(buf[0:])
	if err != nil {
		panic(err)
	}
	fmt.Println(string(buf[:n]))
	time.Sleep(time.Second * 1)
}
