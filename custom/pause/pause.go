package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func main() {
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, syscall.SIGCONT)

	<-channel
	command := os.Args[1]
	args := os.Args[2:]
	cmd := exec.Command(command, args...)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error on running command: %s", err)
	}
}

