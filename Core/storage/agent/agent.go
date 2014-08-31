package agent

import (
	"github.com/joernweissenborn/aursirrt/core/storage"
	"github.com/joernweissenborn/aursirrt/config"
)

func NewAgent(cfg config.RtConfig) (agent StorageAgent) {
	agent.storageReadChannel = make(chan storage.StorageFunc)
	agent.storageWriteChannel = make(chan storage.StorageFunc)

	var sc storage.StorageCore

	sc.Run(cfg, agent.storageWriteChannel,agent.storageReadChannel)

	return
}

type StorageAgent struct {
	storageWriteChannel chan storage.StorageFunc
	storageReadChannel chan storage.StorageFunc
}

func (a StorageAgent) Write(fun storage.StorageFunc) {
	a.storageWriteChannel <- fun
}

func (a StorageAgent) Read(fun storage.StorageFunc) {
	a.storageReadChannel <- fun
}

