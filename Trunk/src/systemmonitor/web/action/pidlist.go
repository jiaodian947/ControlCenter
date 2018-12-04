package action

import (
	"encoding/json"
	"systemmonitor/monitor"

	"github.com/lunny/tango"
)

type Pidlist struct {
	RenderBase
	tango.Json
}

func (p *Pidlist) Get() interface{} {
	ps := monitor.GetAllProcess(100)

	b, err := json.Marshal(ps)
	if err != nil {
		return map[string]interface{}{
			"status": 500,
			"err":    err.Error(),
		}
	}
	return map[string]interface{}{
		"status":  200,
		"sysinfo": string(b),
	}
}
