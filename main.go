package main

import (
	"log"
	"github.com/joernweissenborn/aursirrt/core"
	"github.com/joernweissenborn/aursirrt/dock"
	"os"
	"os/signal"
)

func main(){

	log.Println("AurSirRT launching")

	quit := make(chan struct {})

	aic := make(chan core.AppMessage,100)
	aoc := make(chan core.AppMessage,100)


	core.Launch(aic,aoc)

	dock.Launch(aic,aoc)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(){
		for sig := range c {
			log.Println("ShuttingDOwn",sig)
		}
	}()
	<- quit

}

