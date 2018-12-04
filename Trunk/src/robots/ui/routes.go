package ui

import (
	"robots/ui/actions"

	"github.com/lunny/tango"
)

func setRoutes(t *tango.Tango) {
	t.Get("/", new(actions.Index))
	setting := new(actions.Setting)
	t.Route("POST:PostServer", "/setting/server", setting)
	t.Route("POST:PostRobot", "/setting/robot", setting)
	t.Post("/robot/control", new(actions.Control))
}
