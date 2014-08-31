package callmodel

import "github.com/joernweissenborn/propertygraph2go"

type CallGraph struct {
	*propertygraph2go.InMemoryGraph
}

func NewCallGraph() (callgraph *CallGraph){
	callgraph.Init()
	return
}
