package processor

import (
	"log"
	"storage"
)

const DEBUG = true


func Process(procchan chan Processor, storageagent storage.StorageAgent, maxprocesses int64) {

	print("Initialized")

	procslots := make(chan struct {},maxprocesses)

	for proc := range procchan {
		procslots <- struct{}{}

		proc.Init(procchan, storageagent)

		debugPrint("Processing ")

		go process(proc,procslots)
	}

}

func process(p Processor, ps chan struct{}){
	p.Process()
	<- ps
}

func print(msg string){
	log.Println("PROCESSOR",msg)

}

func debugPrint(msg string){

	if DEBUG {
		log.Println("DEBUG PROCESSOR",msg)

	}

}
