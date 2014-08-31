package callmodel

import (
	"github.com/joernweissenborn/aursir4go/callmodel"
	"github.com/joernweissenborn/aursirrt/storagecore"
)

type CallModelProcessor struct {
	callgraph CallGraph
	storagecoreAgent storagecore.StorageCoreAgent
}

func ProcessCallModel(model callmodel.AurSirCallModel){
	callgraph := NewCallGraph()
	callgraph.CreateVertex("root",nil)
	ProcessFork(callgraph,"root",model.RootFork,model)
}

func ProcessFork(callgraph *CallGraph,
	originId string,fork callmodel.AurSirFork, model callmodel.AurSirCallModel){



}
