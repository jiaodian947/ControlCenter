// +build linux

package action

import "github.com/zcalusic/sysinfo"

func GetSysInfo() SystemInfo {
	var tmp sysinfo.SysInfo
	tmp.GetSysInfo()

	si := SystemInfo{}
	si.OS = tmp.OS.Name
	si.Release = tmp.OS.Release
	si.Kernel = tmp.Kernel.Release
	si.Cpu = []CpuInfo{
		{
			CPU:       0,
			ModelName: tmp.CPU.Model,
			Cores:     int32(tmp.CPU.Cores),
		},
	}
	si.Memory = uint64(tmp.Memory.Size * 1024 * 1024)
	return si
}
