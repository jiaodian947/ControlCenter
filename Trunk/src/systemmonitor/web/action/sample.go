package action

import (
	"errors"
	"systemmonitor/proxy"
	"systemmonitor/sample"

	"github.com/lunny/tango"
)

type Sampler interface {
	Sample(pids string)
	StopSample()
	Query() <-chan *sample.ProcSample
}

type Sample struct {
	RenderBase
	tango.Json
}

func (p *Sample) StopSample() interface{} {
	proxy.Ctx.Sample.StopSample()

	return map[string]interface{}{
		"status": 200,
	}
}

func (p *Sample) QuerySample() interface{} {
	ret := make([]*sample.ProcSample, 0, 16)
	if proxy.Ctx.Sample != nil {
	L:
		for {
			select {
			case s := <-proxy.Ctx.Sample.Query():
				ret = append(ret, s)
			default:
				break L
			}
		}
	}

	return map[string]interface{}{
		"status": 200,
		"data":   ret,
	}
}

func (p *Sample) BeginSample() interface{} {
	type jsondata struct {
		Pids string `json:"pids"`
	}
	var data jsondata
	err := p.DecodeJSON(&data)
	if err != nil {
		return map[string]interface{}{
			"err:":   err,
			"status": 500,
		}
	}

	if proxy.Ctx.Sample != nil {
		proxy.Ctx.Sample.Sample(data.Pids)
		return map[string]interface{}{
			"status": 200,
		}
	}

	return map[string]interface{}{
		"err:":   errors.New("sample not found"),
		"status": 500,
	}
}

func (p *Sample) QueryGameInfo() interface{} {
	return map[string]interface{}{
		"status": 200,
		"data":   proxy.Ctx.Sample.GameSample(),
	}
}
