package infos

import (
	"encoding/json"
	"math"
	"time"

	"github.com/lunny/tango"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

type Infos struct {
	tango.Ctx
	tango.Json
}

type SysInfo struct {
	CpuLoad int    `json:"cpuload"`
	MemLoad int    `json:"memload"`
	NetLoad int    `json:"netload"`
	NetSend uint64 `json:"netsend"`
	NetRecv uint64 `json:"netrecv"`
	oldSend uint64
	oldRecv uint64
}

var (
	old = SysInfo{0, 0, 0, 0, 0, 0, 0}
)

func GetSysInfo() string {
	go func() {
		l, _ := cpu.Percent(time.Second, false)
		old.CpuLoad = int(math.Ceil(l[0]))
	}()
	m, _ := mem.VirtualMemory()
	old.MemLoad = int(math.Ceil(m.UsedPercent))
	n, _ := net.IOCounters(false)
	old.NetSend = n[0].BytesSent - old.oldSend
	old.NetRecv = n[0].BytesRecv - old.oldRecv
	old.oldSend = n[0].BytesSent
	old.oldRecv = n[0].BytesRecv
	b, _ := json.Marshal(old)
	return string(b)
}

func (i *Infos) Get() interface{} {
	h := i.Ctx.Header()
	h.Add("Access-Control-Allow-Origin", "*")
	h.Add("Access-Control-Allow-Methods", "*")
	return map[string]interface{}{
		"status":  200,
		"sysinfo": GetSysInfo(),
	}
}
