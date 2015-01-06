package types

import (
	"storage"
	"github.com/joernweissenborn/aursir4go"
	"log"
)




type Import struct {
	agent storage.StorageAgent
	appid string
	key aursir4go.AppKey
	tags []string
	id string
}


func GetImport(appid string, key aursir4go.AppKey, tags []string, agent storage.StorageAgent) Import {
	i :=  Import{agent,appid,key,tags,""}
	i.setId()
	return i
}

func (i *Import) Exists() bool {
	return i.id == ""

}
func (i *Import) Add() {
	id := make(chan string)
	defer close(id)
	a := GetApp(i.appid, i.agent)
	if !a.Exists() {
		log.Println("STORAGECORE", "Adding Importer failed, app does not exist:", i.appid)
		return
	}
	k := GetAppKey(i.key, i.agent)
	k.Create()
	keyid := k.GetId()
	i.agent.Write(func(sc *storage.StorageCore) {
		av := sc.InMemoryGraph.GetVertex(a.Id)
		kv := sc.InMemoryGraph.GetVertex(keyid)
		ev := sc.InMemoryGraph.CreateVertex(storage.GenerateUuid(), nil)


		sc.InMemoryGraph.CreateEdge(storage.GenerateUuid(), storage.IMPORT_EDGE, kv, ev, nil)
		sc.InMemoryGraph.CreateEdge(storage.GenerateUuid(), storage.IMPORT_EDGE, ev, av, nil)

		id <- ev.Id
	})
	i.id = <-id

	if i.id != "" {
		for _, tag := range i.tags {
			t := GetTag(k,tag,i.agent)
			t.Create()

		}
	}
}

func (e *Import) GetApp() App{
	return GetApp(e.appid,e.agent)
}
func (e *Import) GetAppKey() AppKey{
	return GetAppKey(e.key,e.agent)
}
func (e *Import) setId() {
	keyid := e.GetAppKey().GetId()
	a := GetApp(e.appid, e.agent)
	if !a.Exists() {
		log.Println("STORAGECORE", "Setting Importerid failed, app does not exist:", e.appid)
		return
	}
	c := make(chan string)
	defer close(c)
	e.agent.Read(func (sc *storage.StorageCore){
		av := sc.GetVertex(e.appid)
		i := 0
		//app - ImportEDGE > Import
		for _,Importedge := range av.Outgoing{
			if Importedge.Label == storage.IMPORT_EDGE {
				//Import - ImportEDGE > Key
				Import := Importedge.Head
				for _,Importkeyedge := range Import.Outgoing {
					if Importkeyedge.Label == storage.IMPORT_EDGE {
						if keyid == Importkeyedge.Head.Id {

							//Import - TAGEDGE > Tag
							for _, tagedge := range Import.Outgoing {
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
								c <- Import.Id
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
func (e Import) GetId() string {
	return e.id
}

func (e Import) UpdateTags(){
	e.ClearTags()
	e.agent.Write(func (sc *storage.StorageCore) {

	})
}

func (e Import) ClearTags(){
	//key := GetAppKey(e.key,e.agent)

}

func (e Import) GetTags() ([]Tag){
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
