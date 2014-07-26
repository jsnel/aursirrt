package main

import (
	"log"
	"github.com/joernweissenborn/AurSirRt/Dock"
	"github.com/joernweissenborn/AurSirRt/Core"
)

func main(){

	log.Println("AurSirRT launching")

	quit := make(chan struct {})

	aic := make(chan Core.AppMessage,100)
	aoc := make(chan Core.AppMessage,100)


	Core.Launch(aic,aoc)

	Dock.Launch(aic,aoc)

	<- quit

}

