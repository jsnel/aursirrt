package boot

import (
	"log"
	"aursirrt/src/processor"
	"aursirrt/src/storage"
	"aursirrt/src/dock"
	"aursirrt/src/dock/dockzmq"
	"aursirrt/src/cmdlineinterface"
	"flag"
	"aursirrt/src/config"
)

const (
	MAX_PROCESSORS = 4
)

func Boot(){

	mprint("AurSir RT starting")
	flag.Parse()

	a := bootStorage()
	                   id:= a.GetId()
	p := bootCore(a)
	var z dockzmq.DockerZmq
	z.SetPort(int64(*config.Zmqport))
	bootDocker(p, z, id)

//	var w dockwebsockets.DockerWebSockets
//	bootDocker(p,w)

	bootCmdlineinterface(p)
}
func BootWithoutCmdlineinterface(){

	mprint("AurSir RT starting")
	flag.Parse()
	a := bootStorage()

	id:= a.GetId()
	p := bootCore(a)
	var z dockzmq.DockerZmq
	z.SetPort(int64(*config.Zmqport))

	bootDocker(p, z, id)

	//	var w dockwebsockets.DockerWebSockets
//	bootDocker(p,w)
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

func bootDocker(p chan processor.Processor, d dock.Docker, id string) {
	mprint("Launching Dock")
	agent := dock.NewAgent(p)
	d.Launch(agent,id)

}


func bootCmdlineinterface(p chan processor.Processor) {
	mprint("Launching Cmdlineinterface")
	cli := cmdlineinterface.CmdLineInterface{}
	cli.Run()

}


func mprint(msg string){

	log.Println("BOOT",msg)

}
