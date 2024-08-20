package main

import (
	"os"
	"os/signal"
	"sao/discord"
	"sao/world"
	"syscall"
)

func main() {
	world := world.CreateWorld()

	world.LoadBackup()

	go world.StartClock()

	discord.World = &world

	go discord.StartClient()

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
