package proxy

import (
	"systemmonitor/sample"
)

type Context struct {
	Sample *sample.Sample
}

var (
	Ctx Context
)
