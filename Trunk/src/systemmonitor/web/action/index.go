package action

import (
	"fmt"

	"github.com/tango-contrib/renders"
)

type Index struct {
	RenderBase
}

type CpuInfo struct {
	CPU       int32
	ModelName string
	Cores     int32
}

type SystemInfo struct {
	OS      string
	Release string
	Kernel  string
	Cpu     []CpuInfo
	Memory  uint64
}

func (i *Index) Get() error {
	si := GetSysInfo()
	fmt.Println(si)
	return i.Render("index.html", renders.T{
		"system": si,
	})
}
