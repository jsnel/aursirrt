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
	"net"
)

const (
	MAX_PROCESSORS = 8
)

func Boot(){

	mprint("AurSir RT starting")
	flag.Parse()

	a := bootStorage()
	                   id:= a.GetId()
	mprint("Nodeid is "+id)

	p := bootCore(a)
	if *config.P2p {
		var nz dockzmq.DockerZmq
		nz.SetIp(*config.Zmqip)
		nz.SetP2P(true)
		nz.SetPort(getRandomPort())
		bootDocker(p, nz, id)
	}
	var lz dockzmq.DockerZmq
	lz.SetPort(int64(*config.Zmqport))
	bootDocker(p, lz, id)

//	var w dockwebsockets.DockerWebSockets
//	bootDocker(p,w)

	bootCmdlineinterface(p)
}
func BootWithoutCmdlineinterface(){

	mprint("AurSir RT starting")
	flag.Parse()
	a := bootStorage()

	id:= a.GetId()
	mprint("Nodeid is "+id)
	p := bootCore(a)

	var z dockzmq.DockerZmq
	z.SetPort(int64(*config.Zmqport))
	z.SetIp(*config.Zmqip)
	z.SetP2P(*config.P2p)
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


func getRandomPort() int64 {
	l, err := net.Listen("tcp", "127.0.0.1:0") // listen on localhost
	if err != nil {
		panic("Could not find a free port")
	}
	defer l.Close()
	return int64(l.Addr().(*net.TCPAddr).Port)
}
