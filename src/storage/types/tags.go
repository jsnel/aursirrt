package types

import (
	"storage"
)

type Tag struct {
	agent storage.StorageAgent
	key AppKey
	name string
	id string
}

func GetTag(key AppKey, name string, agent storage.StorageAgent)Tag{
	tag := Tag{agent,key,name,""}
	tag.setId()
	return tag
}

func (t *Tag) setId() {
	c := make(chan string)
	defer close(c)
	keyid := t.key.GetId()
	t.agent.Read(func (sc *storage.StorageCore){
		kv := sc.GetVertex(keyid)
		if kv == nil {
			c <-""
			return
		}
		for _,e := range kv.Incoming{
			if e.Label == storage.TAG_EDGE {
				tagname,_ := e.Tail.Properties.(string)
				if tagname == t.name {
					c<-e.Tail.Id
					return
				}
			}
		}
		c <- ""
	})
	t.id = <- c
}
func (t Tag) GetId() string{
	if t.id == ""{
		t.setId()
	}
	if !t.Exists() {
		return ""
	}

	return t.id



}

func (t Tag) Create(){
	if !t.Exists() {
		c := make(chan string)
		defer close(c)
		keyid := t.key.GetId()
		t.agent.Write(func (sc *storage.StorageCore){
			tv := sc.CreateVertex(storage.GenerateUuid(), t.name)
			kv := sc.GetVertex(keyid)
			sc.CreateEdge(storage.GenerateUuid(), storage.TAG_EDGE, kv,tv, nil)

			c <- tv.Id
		})
		t.id = <-c
	}
}



func (t Tag) LinkExport(e *Export){
	if t.Exists() {
		c := make(chan string)
		defer close(c)
		eid := e.GetId()

		t.agent.Write(func (sc *storage.StorageCore){
			tv := sc.GetVertex(t.GetId())
			ev := sc.GetVertex(eid)
			sc.CreateEdge(storage.GenerateUuid(), storage.TAG_EDGE, tv, ev, nil)

			c <- ""
		})
	<-c
	}
}

func (t Tag) Exists() bool {
	c := make(chan bool)
	defer close(c)
	keyid := t.key.GetId()
	t.agent.Read(func (sc *storage.StorageCore){
		kv := sc.GetVertex(keyid)
		if kv == nil {
			c <-false
			return
		}
		for _,e := range kv.Incoming{
			if e.Label == storage.TAG_EDGE {
				tagname,_ := e.Tail.Properties.(string)
				if tagname == t.name {
					c<-true
					return
				}
			}
		}
		c <- false
	})
	return <- c
}
