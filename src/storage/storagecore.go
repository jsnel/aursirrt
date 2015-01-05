package storage

import (
	"github.com/joernweissenborn/propertygraph2go"
	uuid "github.com/nu7hatch/gouuid"

	"log"
)

type StorageFunc func(storageCore *StorageCore)


type StorageCore struct {
	*propertygraph2go.InMemoryGraph
	root propertygraph2go.Vertex
}

func (sc *StorageCore) ExecuteFunc(storFunc StorageFunc){
	storFunc(sc)
}

func (sc *StorageCore) Run(storageWriteChan,storageReadChan chan StorageFunc){
	sc.InMemoryGraph = propertygraph2go.New()

	sc.root = sc.InMemoryGraph.CreateVertex("root",nil})

	ok := true

	for ok {
		select {

		case fun, ok := <-storageWriteChan:
			if ok{
				sc.ExecuteFunc(fun)
			}

		case fun, ok := <-storageReadChan:
			if ok{
				sc.ExecuteFunc(fun)
			}
		}
	}
}
func GenerateUuid() string {
	Uuid, err := uuid.NewV4()
	if err != nil {
		log.Fatal("Failed to generate UUID")
		return ""
	}
	return Uuid.String()
}