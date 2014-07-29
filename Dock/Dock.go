package dock

import (
	"log"
	"github.com/joernweissenborn/aursirrt/core"
)

type registerDockedApp struct{
	AppId string
	AppChan chan core.AppMessage
}

func Launch(ic, oc chan core.AppMessage){

	log.Println("Dock Launching")

	rc := make(chan registerDockedApp)

	for i,docker := range CfgDocker(){

		log.Println("Dock Launching Docker", i+1)

		docker.Launch(ic, rc)
	}

	dockRouter(oc,rc)

}

func dockRouter(mc chan core.AppMessage, rc chan registerDockedApp ){
	routeTable := make(map[string]chan core.AppMessage)
	for{
	select {
	case appregister, ok := <- rc:
		if ok{
			routeTable[appregister.AppId] = appregister.AppChan
			log.Println("Registered out channel for",appregister.AppId)
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
