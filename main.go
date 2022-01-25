package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

// build is the git version of this program. It is set using build flags in the makefile.
var build = "develop"

func main() {
	log.Println("starting service", build)
	defer log.Println("service ended")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	log.Println("stopping service")
}
