package app

import (
	"github.com/joernweissenborn/aursir4go"
	"log"
	"storage"
)


type App struct {
	Agent storage.StorageAgent

Id string
}

func Get(Id string, Agent storage.StorageAgent ) App {
	return App{Agent,Id}
}

func (app App) Exists() bool {

	c := make(chan bool)

	app.Agent.Read(func (sc *storage.StorageCore){
		c <- sc.GetVertex(app.Id) != nil
	})

	return <- c
}

func (app App) Create(DockMessage aursir4go.AurSirDockMessage){
	c := make(chan struct{})

	app.Agent.Write(func (sc *storage.StorageCore){
		log.Println("STORAGECORE","Assinging ID to App",DockMessage.AppName,"->", app.Id)
		sc.CreateVertex(app.Id, DockMessage)
		c <- struct{}{}
	})

	<-c
}
