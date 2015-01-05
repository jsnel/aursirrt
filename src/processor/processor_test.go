package processor

import (
	"testing"
)

type testprocessor struct {
	*GenericProcessor
	c chan bool
}

func (tp testprocessor) Process(){
	tp.c <- true
}

func TestBootProcessingCore(t *testing.T){
	var tp testprocessor
	tp.GenericProcessor = GetGenericProcessor()
	tp.c = make(chan bool)
	pc := make(chan Processor)
	go Process(pc, 1)
	defer close(pc)
	pc <- &tp

	<-tp.c


}

type spawntestprocessor struct {
	*GenericProcessor
	c chan bool
}

func (tsp spawntestprocessor) Process(){
	var tp testprocessor
	tp.GenericProcessor = GetGenericProcessor()
	tp.c = tsp.c
	tsp.SpawnProcess(tp)
}

func TestSpawnProcess(t *testing.T){
	var tp spawntestprocessor
	tp.GenericProcessor = GetGenericProcessor()
	tp.c = make(chan bool)
	pc := make(chan Processor)
	go Process(pc, 1)
	defer close(pc)
	pc <- &tp

	<-tp.c


}
