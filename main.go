package main

import (
	"log"
	"os"

	"github.com/TheRealShek/mini-docker/runtime"
)

func main() {
	// If orchestrator is present
	if os.Getenv("CONTAINER_INIT") == "1" {
		if err := runtime.ContainerInit(); err != nil {
			log.Fatal(err)
		}
		return
	}
	// orchestrator is not present, so make a Parent Process
	// child detects CONTAINER_INIT and runs setup
	runtime.Run()
}
