package main

import (
	"log"

	"github.com/ginie/sasrd"
)

func main() {
	conf, err := sasrd.NewConfigFromFile("sasrd.json")
	if err != nil {
		log.Printf("[ERROR] %s", err)
	}

	sasrd.StartServer(conf)
}
