package ui

import (
	"html/template"

	"github.com/lunny/tango"
	"github.com/tango-contrib/events"
	"github.com/tango-contrib/renders"
)

type NullLog struct {
}

func (l *NullLog) Debugf(format string, v ...interface{}) {}
func (l *NullLog) Debug(v ...interface{})                 {}
func (l *NullLog) Infof(format string, v ...interface{})  {}
func (l *NullLog) Info(v ...interface{})                  {}
func (l *NullLog) Warnf(format string, v ...interface{})  {}
func (l *NullLog) Warn(v ...interface{})                  {}
func (l *NullLog) Errorf(format string, v ...interface{}) {}
func (l *NullLog) Error(v ...interface{})                 {}

func Serv(port int) {
	t := tango.Classic(&NullLog{})
	//t := tango.Classic()
	t.Use(
		events.Events(),
		tango.Static(tango.StaticOptions{
			RootPath: "./views/statics/assets",
			Prefix:   "assets",
		}),
		renders.New(renders.Options{
			Reload:    true,
			Directory: "./views/templates",
			Funcs:     template.FuncMap{},
		}),
	)
	setRoutes(t)
	t.Run(port)
}
