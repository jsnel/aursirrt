package types

import (
	"log"
	"storage"
	"github.com/joernweissenborn/aursir4go/messages"
	"dock"
)

type appproperties struct {
	dockmsg messages.DockMessage
	connection dock.Connection
}

type App struct {
	agent storage.StorageAgent
	Id string
}

func GetApp(Id string, Agent storage.StorageAgent) App {
	return App{Agent,Id}
}

func (app App) Exists() bool {

	c := make(chan bool)
	defer close(c)
	app.agent.Read(func (sc *storage.StorageCore){
		c <- sc.GetVertex(app.Id) != nil
	})

	return <- c
}

func (app App) GetConnection() dock.Connection {

	c := make(chan dock.Connection)
	defer close(c)
	app.agent.Read(func (sc *storage.StorageCore){
		c <- sc.GetVertex(app.Id).Properties.(appproperties).connection
	})

	return <- c
}

func (app App) Create(DockMessage messages.DockMessage, Connection dock.Connection) bool{
	if app.Exists() {
		return false
	}
	c := make(chan struct{})
	defer close(c)
	app.agent.Write(func (sc *storage.StorageCore){
		log.Println("STORAGECORE","Assinging ID to App",DockMessage.AppName,"->", app.Id)
		sc.CreateVertex(app.Id, appproperties{DockMessage,Connection})
		c <- struct{}{}
	})

	<-c
	return true
}


func (app App) Remove(){
	c := make(chan struct{})
	if !app.Exists() {
		log.Println("STORAGECORE", "Error removing app %d, does not exist", app.Id)
		return
	}
	app.agent.Write(func (sc *storage.StorageCore){
		log.Println("STORAGECORE","Removing App",app.Id)
		sc.RemoveVertex(app.Id)
		c <- struct{}{}
	})

	<-c
}


