package storage

import (
	"github.com/joernweissenborn/propertygraph2go"
	"github.com/joernweissenborn/aursirrt/config"
)

type StorageFunc func(storageCore *StorageCore)


type StorageCore struct {
	propertygraph2go.SemiPersistentGraph
}

func (sc *StorageCore) ExecuteFunc(storFunc StorageFunc){
	storFunc(sc)
}

func (sc *StorageCore) Run(cfg config.RtConfig, storageWriteChan,storageReadChan chan StorageFunc){

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
