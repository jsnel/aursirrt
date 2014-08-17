package dock

import (
	"log"
	"github.com/joernweissenborn/aursirrt/core"
)

type registerDockedApp struct{
	AppId string
	AppChan chan core.AppMessage
}

type ungisterDockedApp struct{
	AppId string
}

func Launch(ic, oc chan core.AppMessage){

	log.Println("Dock Launching")

	rc := make(chan interface {})


	for i,docker := range CfgDocker(){

		log.Println("Dock Launching Docker", i+1)

		docker.Launch(ic, rc)
	}

	dockRouter(oc,rc)

}

func dockRouter(mc chan core.AppMessage, rc chan interface {} ){
	routeTable := make(map[string]chan core.AppMessage)
	for{
	select {
	case msg, ok := <- rc:
		if ok{
			switch req := msg.(type){
			case registerDockedApp:
				routeTable[req.AppId] = req.AppChan
				log.Println("Registered out channel for",req.AppId)
			case ungisterDockedApp:
				close(routeTable[req.AppId])
				delete(routeTable,req.AppId)
			}
		}

	case appmsg, ok := <- mc:
		if ok {
			ac, f := routeTable[appmsg.SenderUUID]
			if f {
				ac <- appmsg
			}
		}
	}
	}
}
