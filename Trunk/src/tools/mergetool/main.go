package main

import (
	"log"
	"tools/mergetool/merge"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	config := merge.NewConfig()
	if err := config.Load("merge.xml"); err != nil {
		panic(err)
	}

	log.Println(config)
	m := merge.New(config)
	err := m.Main()
	if err != nil {
		log.Println(err)
	}
	m.Exit()
}
