// +build windows

package action

import (
	"github.com/matishsiao/goInfo"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func GetSysInfo() SystemInfo {
	si := SystemInfo{}

	os := goInfo.GetInfo()
	si.OS = os.OS
	si.Kernel = os.Kernel
	si.Release = os.Core

	info, err := cpu.Info()
	if err == nil {
		si.Cpu = make([]CpuInfo, 0, len(info))
		for _, v := range info {
			ci := CpuInfo{}
			ci.CPU = v.CPU
			ci.ModelName = v.ModelName
			ci.Cores = v.Cores
			si.Cpu = append(si.Cpu, ci)
		}

	}

	memory, err := mem.VirtualMemory()
	if err == nil {
		si.Memory = memory.Total
	}

	return si
}
