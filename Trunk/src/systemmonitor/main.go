package main

import (
	"systemmonitor/proxy"
	"systemmonitor/sample"
	"systemmonitor/web"

	"github.com/mysll/toolkit"
)

func main() {
	s := sample.NewSample(128, "192.168.1.134", 28010)
	proxy.Ctx = proxy.Context{Sample: s}
	go web.Serv(9091)
	toolkit.WaitForQuit()
	return
}
