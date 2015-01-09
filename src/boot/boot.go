package boot

import (
	"log"
	"processor"
	"storage"
	"dock"
	"dock/dockzmq"
	"dock/dockwebsockets"
)

const (
	MAX_PROCESSORS = 4
)

func Boot(){

	mprint("AurSir RT starting")

	a := bootStorage()

	p := bootCore(a)
	var z dockzmq.DockerZmq
	bootDocker(p, z)

	var w dockwebsockets.DockerWebSockets
	bootDocker(p,w)
}

func bootStorage() storage.StorageAgent {
	mprint("Launching Storage")

	return storage.NewAgent()
}


func bootCore(a storage.StorageAgent) (processingChan chan processor.Processor){

	mprint("Launching Core")

	processingChan = make(chan processor.Processor)

	go processor.Process(processingChan, a, MAX_PROCESSORS)

	return
}

func bootDocker(p chan processor.Processor, d dock.Docker) storage.StorageAgent {
	mprint("Launching Dock")
	agent := dock.DockAgent{p}
	d.Launch(agent)

	return storage.NewAgent()
}


func mprint(msg string){

	log.Println("BOOT",msg)

}
