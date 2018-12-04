package monitor

import (
	"github.com/shirou/gopsutil/process"
)

type Process struct {
	Cmd string
	Pid int32
}

func GetAllProcess(max int) []Process {
	p, err := process.Processes()
	if err != nil {
		return nil
	}

	ret := make([]Process, 0, len(p))
	for _, ps := range p {
		cl, err := ps.Cmdline()
		if err != nil || cl == "" {
			continue
		}
		ret = append(ret, Process{Cmd: cl, Pid: ps.Pid})
	}
	return ret
}
