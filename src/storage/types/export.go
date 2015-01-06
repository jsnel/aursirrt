package types

import (
	"storage"
	"github.com/joernweissenborn/aursir4go"
	"log"
)




type Export struct {
	agent storage.StorageAgent
	appid string
	key aursir4go.AppKey
	tags []string
	id string
}


func GetExport(appid string, key aursir4go.AppKey, tags []string, agent storage.StorageAgent) Export {
	e :=  Export{agent,appid,key,tags,""}
	return e
}

func (e *Export) Exists() bool {
	return e.id == ""

}
func (e *Export) Add() {
	id := make(chan string)
	defer close(id)
	a := GetApp(e.appid, e.agent)
	if !a.Exists() {
		log.Println("STORAGECORE", "Adding exporter failed, app does not exist:", e.appid)
		return
	}
	k := GetAppKey(e.key, e.agent)
	k.Create()
	keyid := k.GetId()
	e.agent.Write(func(sc *storage.StorageCore) {
		av := sc.InMemoryGraph.GetVertex(a.Id)
		kv := sc.InMemoryGraph.GetVertex(keyid)
		ev := sc.InMemoryGraph.CreateVertex(storage.GenerateUuid(), nil)


		sc.InMemoryGraph.CreateEdge(storage.GenerateUuid(), storage.EXPORT_EDGE, kv, ev, nil)
		sc.InMemoryGraph.CreateEdge(storage.GenerateUuid(), storage.EXPORT_EDGE, ev, av, nil)

		id <- ev.Id
	})
	e.id = <-id

	if e.id != "" {
		for _, tag := range e.tags {

			t := GetTag(k,tag,e.agent)

			t.Create()
			t.LinkExport(e)
			log.Println("STORAGECORE", "Adding exporter failed, app does not exist:", t.Exists())

		}
	}
}

func (e *Export) GetApp() App{
	return GetApp(e.appid,e.agent)
}
func (e *Export) GetAppKey() AppKey{
	return GetAppKey(e.key,e.agent)
}
func (e *Export) setId() {
	keyid := e.GetAppKey().GetId()
	a := GetApp(e.appid, e.agent)
	if !a.Exists() {
		log.Println("STORAGECORE", "Setting exporterid failed, app does not exist:", e.appid)
		return
	}
	c := make(chan string)
	defer close(c)
	e.agent.Read(func (sc *storage.StorageCore){
		av := sc.GetVertex(e.appid)

		i := 0
		//app - EXPORTEDGE > Export
		for _,exportedge := range av.Outgoing{
			if exportedge.Label == storage.EXPORT_EDGE {

				//Export - EXPORTEDGE > Key
				export := exportedge.Head
				for _,exportkeyedge := range export.Outgoing {
					if exportkeyedge.Label == storage.EXPORT_EDGE {
						if keyid == exportkeyedge.Head.Id {
							log.Println("STORAGECORE",len(export.Outgoing))

							//Export - TAGEDGE > Tag
							for _, tagedge := range export.Outgoing {
								if tagedge.Label == storage.TAG_EDGE {
									tagname := tagedge.Head.Properties.(string)
									for _, tn := range e.tags {

										if tn == tagname {
											i++
											break
										}
									}
								}
							}

							if len(e.tags) == i {
								c <- export.Id
								return
							}

						}
					}
				}
			}
		}
		c  <- ""
	})
	e.id = <- c
	return
}
func (e *Export) GetId() string {
	if e.id == "" {
		e.setId()
	}
	return e.id
}

func (e Export) UpdateTags(){
	e.ClearTags()
	e.agent.Write(func (sc *storage.StorageCore) {

	})
}

func (e Export) ClearTags(){
	//key := GetAppKey(e.key,e.agent)

}

func (e Export) GetTags() ([]Tag){
	tags := []Tag{}
	if e.GetId() == "" {
		return tags
	}
	k := GetAppKey(e.key, e.agent)
	c := make (chan []Tag)
	defer close(c)
	e.agent.Read(func (sc *storage.StorageCore){
		ev := sc.GetVertex(e.GetId())
		for _,tagedge := range ev.Outgoing{
			if tagedge.Label == storage.TAG_EDGE {
				tagname,_ := tagedge.Head.Properties.(string)
				tags = append(tags,Tag{e.agent,k,tagname,tagedge.Head.Id})
			}
		}
		c <- tags
	})
	return <-c
}
