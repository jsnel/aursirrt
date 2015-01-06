package types

import (
	"storage"
	"github.com/joernweissenborn/aursir4go"
	"log"
)

type AppKey struct {
	agent storage.StorageAgent
	appkey aursir4go.AppKey
}

func GetAppKey(appkey aursir4go.AppKey, agent storage.StorageAgent) AppKey{
	return AppKey{agent,appkey}
}

func (a AppKey) Exists() bool {

	c := make(chan bool)

	a.agent.Read(func (sc *storage.StorageCore){
		for _, ke := range sc.Root.Outgoing {
			if ke.Label == storage.KNOWN_APPKEY_EDGE {
				key, _ := ke.Head.Properties.(aursir4go.AppKey)
				if key.ApplicationKeyName == a.appkey.ApplicationKeyName {
					c<-true
					return
				}
			}
		}
		c<-false
	})

	return <- c
}

func (a AppKey) GetId() string {

	c := make(chan string)

	a.agent.Read(func (sc *storage.StorageCore){
		for _, ke := range sc.Root.Outgoing {
			if ke.Label == storage.KNOWN_APPKEY_EDGE {
				key, _ := ke.Head.Properties.(aursir4go.AppKey)
				if key.ApplicationKeyName == a.appkey.ApplicationKeyName {
					c<-ke.Head.Id
					return
				}
			}
		}
		c<-""
	})

	return <- c
}

func (a AppKey) Create() {
	c := make(chan bool)
	if !a.Exists() {
		a.agent.Write(func(sc *storage.StorageCore) {
			kv := sc.InMemoryGraph.CreateVertex(storage.GenerateUuid(), a.appkey)

			sc.InMemoryGraph.CreateEdge(storage.GenerateUuid(), storage.KNOWN_APPKEY_EDGE, kv, sc.Root, nil)

			log.Println("StorageCore Key registered")
			c<-false
		})
		<-c
	}
	return
}
