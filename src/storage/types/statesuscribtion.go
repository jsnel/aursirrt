package types

import (
	"storage"
	"dock/connection"
)

type StateSuscribtion struct {
	agent storage.StorageAgent
	connection connection.Connection
}

func GetAllSuscribtions(agent storage.StorageAgent) []StateSuscribtion {
	c:= make(chan StateSuscribtion)
	agent.Read(func (sc *storage.StorageCore){
		for _, suscribeedge := range sc.Root.Outgoing {
			if suscribeedge.Label == STATE_SUSCRIBTION_EDGE {
				c <- StateSuscribtion{agent, suscribeedge.Head.Properties.(connection.Connection)}
			}
		}
		close(c)
	})
	suscribtions := []StateSuscribtion{}
	for suscribtion := range c {
		suscribtions = append(suscribtions,suscribtion)
	}
	return suscribtions
}

func CreateSuscribtion(connection connection.Connection,agent storage.StorageAgent) {
	suscribeid := storage.GenerateUuid()
	agent.Read(func(sc *storage.StorageCore) {
	   sv := sc.CreateVertex(suscribeid,connection)
		sc.CreateEdge(storage.GenerateUuid(), STATE_SUSCRIBTION_EDGE,sc.Root,sv,nil)
	})
}
