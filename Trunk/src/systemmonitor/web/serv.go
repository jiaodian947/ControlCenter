package web

import (
	"fmt"
	"html/template"
	"math"

	"github.com/lunny/tango"
	"github.com/tango-contrib/events"
	"github.com/tango-contrib/renders"
)

func memTotal(bytes uint64) string {
	if bytes == 0 {
		return "0B"
	}

	k := float64(1024)
	sizes := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
	i := int(math.Floor(math.Log(float64(bytes)) / math.Log(k)))

	return fmt.Sprintf("%.0f%s", (float64(bytes) / math.Pow(k, float64(i))), sizes[i])
}

func Serv(port int) {
	t := tango.Classic()
	t.Use(
		events.Events(),
		tango.Static(tango.StaticOptions{
			RootPath: "./views/assets",
			Prefix:   "assets",
		}),
		renders.New(renders.Options{
			Reload:    true,
			Directory: "./views/templates",
			Funcs: template.FuncMap{
				"MemTotal": memTotal,
			},
		}),
	)
	setRoutes(t)
	t.Run(port)
}
