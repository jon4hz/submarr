package main

import (
	"log"

	"github.com/jon4hz/submarr/cmd"
	zone "github.com/lrstanley/bubblezone"
)

func main() {
	zone.NewGlobal()

	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
