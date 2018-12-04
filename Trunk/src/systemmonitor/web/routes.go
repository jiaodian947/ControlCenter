package web

import (
	"systemmonitor/web/action"

	"github.com/lunny/tango"
)

func setRoutes(t *tango.Tango) {
	t.Get("/", new(action.Index))
	t.Get("/pid", new(action.Pidlist))
	s := new(action.Sample)
	t.Route("POST:BeginSample", "/sample", s)
	t.Route("GET:StopSample", "/stop", s)
	t.Route("GET:QuerySample", "/query", s)
	t.Route("GET:QueryGameInfo", "/info", s)
}
