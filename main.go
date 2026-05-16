package main

import (
	"log"

	"github.com/TheRealShek/mini-docker/runtime"
)

func main() {
	if err := runtime.Run(); err != nil {
		log.Fatal(err)
	}
}
