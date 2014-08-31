package types

import (
	"github.com/joernweissenborn/aursirrt/core/storage"
	"github.com/joernweissenborn/aursir4go"
	"log"
)


type App struct {
	StorageType
	Id string
}

func (app App) Exists(Id string) bool {

	c := make(chan bool)

	app.Agent.Read(func (sc *storage.StorageCore){
		c <- sc.GetVertex(Id) != nil
	})

	return <- c
}

func (app App) Create(Id string, DockMessage aursir4go.AurSirDockMessage){
	c := make(chan struct{})

	app.Agent.Write(func (sc *storage.StorageCore){
		log.Println("STORAGECORE","Assinging ID to App",DockMessage.AppName,"->", Id)
		sc.CreateVertex(Id, DockMessage)
		c <- struct{}{}
	})

	<-c
}
