package export

import (
	"storage"
	"github.com/joernweissenborn/aursir4go"
	"log"
	"storage/types/app"
	"storage/types/appkey"
)




type Export struct {
	agent storage.StorageAgent
	appid string
	key aursir4go.AppKey
	tags []string
	id string
}




func Get(appid string, key aursir4go.AppKey, tags []string, agent storage.StorageAgent) Export {
	return Export{agent,appid,key,tags,""}
}

func (e Export) Add() {
	id := make (chan string)
	a := app.Get(e.appid,e.agent)
	if !a.Exists() {
		log.Println("StorageCore adding exporter failed, app does not exist:", e.appid)
		id <- ""
		return
	}
	k := appkey.Get(e.key,e.agent)
	k.Create()
	keyid := k.GetId()
	e.agent.Write(func (sc *StorageCore) {
		av := sc.InMemoryGraph.GetVertex(a.Id)
		kv := sc.InMemoryGraph.GetVertex(keyid)
		ev := sc.InMemoryGraph.CreateVertex(storage.GenerateUuid(), nil)


		sc.graph.CreateEdge(storage.GenerateUuid(), storage.EXPORT_EDGE, kv, ev, nil)
		sc.graph.CreateEdge(storage.GenerateUuid(), storage.EXPORT_EDGE, ev, av, nil)

		id <- ev.Id}
	})
	e.id =  <- id

	if e.id != "" {

	}

func (e Export) UpdateTags(){
	e.ClearTags()
	e.agent.Write(func (sc *StorageCore) {

	}
	}

func (e Export) ClearTags(){
	key := appkey.Get(e.key,e.agent)

}
