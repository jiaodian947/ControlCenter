package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

var (
	addr = flag.String("h", "", "host")
	port = flag.Int("p", 7888, "port")
)

func main() {
	flag.Parse()

	exitChan := make(chan int)
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		exitChan <- 1
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *addr, *port))
	if err != nil {
		panic(err)
	}
	log.Println("start at:", *addr, *port)
	go serv(l)
	<-exitChan
	l.Close()
	log.Println("quit")
}

func serv(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				log.Printf("NOTICE: temporary Accept() failure - %s", err)
				runtime.Gosched()
				continue
			}
			if !strings.Contains(err.Error(), "use of closed network connection") {
				log.Printf("ERROR: listener.Accept() - %s", err)
			}
			break
		}

		go handler(conn)

	}
}

func handler(conn net.Conn) {

	r := bufio.NewReaderSize(conn, 4096)
	w := bufio.NewWriterSize(conn, 4096)
	w.WriteString("welcome to build console\r\n")
	w.Flush()
	for {
		w.WriteString(">")
		w.Flush()
		l, err := r.ReadSlice('\n')
		if err != nil {
			break
		}

		line := strings.TrimSpace(string(l))
		line = strings.Replace(strings.Replace(line, "\r", "", -1), "\n", "", -1)
		params := strings.Split(line, " ")
		if len(params) > 0 {
			if params[0] == "quit" {
				break
			}
			execCmd(w, params)
		}
	}
	w.WriteString("bye")
	w.Flush()
	conn.Close()
	log.Println("client quit")
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func execCmd(w *bufio.Writer, params []string) {
	switch params[0] {
	case "build":
		cmd := getCurrentDirectory() + "/build.bat"
		w.WriteString("Run " + cmd + "\r\n")
		exec := exec.Command("cmd.exe", "/c", cmd)

		stdoutPipe, err := exec.StdoutPipe()
		if err != nil {
			panic(err)
		}
		exec.Start()
		go io.Copy(w, stdoutPipe)
		exec.Wait()
		w.WriteString(" build complete\r\n")
		w.Flush()
	case "help":
		w.WriteString(" USAGE:\thelp: show help\r\n\tbuild: build server\r\n\tquit: quit console\r\n")
		w.Flush()
	default:
		w.WriteString(" unknown command\r\n")
		w.Flush()
	}
}
