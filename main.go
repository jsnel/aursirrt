package main

import (
	"github.com/joernweissenborn/aursirrt/core"
	"github.com/joernweissenborn/aursirrt/dock"
	"log"
	"os"
	"os/signal"

	"github.com/joernweissenborn/aursirrt/config"
	"github.com/joernweissenborn/aursirrt/webserver"
)

func main() {
	log.Println("AurSirRT launching")

	log.Println("AurSirRT loading config")

	cfg, err := config.NewRtCfg()

	if err != nil {
		log.Fatal("AurSirRT error loding cfg", err)
	}
	log.Println("AurSirRT config loaded")

	quit := make(chan struct{})

	aic := make(chan core.AppMessage, 100)
	//aic1 := make(chan core.AppMessage,100)
	//aoc1 := make(chan core.AppMessage,100)
	aoc := make(chan core.AppMessage, 100)

	//go debug("AIC",aic,aic1)
	//go debug("AOC",aoc,aoc1)

	go webserver.Launch(cfg)

	core.Launch(aic, aoc, cfg)

	dock.Launch(aic, aoc)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Println("ShuttingDOwn", sig)
		}
	}()
	<-quit

}

func debug(name string, in, out chan core.AppMessage) {
	for msg := range in {
		log.Println("DEBUG", name, msg)
		out <- msg
	}
}
