package main

import (
	"log"

	"github.com/jon4hz/subrr/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
