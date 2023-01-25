package main

import (
	"log"

	"github.com/jon4hz/subrr/cmd"
	zone "github.com/lrstanley/bubblezone"
)

func main() {
	zone.NewGlobal()

	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
