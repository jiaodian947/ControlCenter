package main

import (
	"robots/deamon/infos"

	"github.com/lunny/tango"
)

func main() {
	t := tango.Classic()
	t.Get("/sysinfo", new(infos.Infos))
	t.Run(8011)
}
