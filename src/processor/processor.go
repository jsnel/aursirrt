package processor

import "storage/agent"

type Processor interface {
	Process()
	Init(chan Processor, agent.StorageAgent)
	SpawnProcess(Processor)
}

func GetGenericProcessor() (gp *GenericProcessor) {
	var proc GenericProcessor
	gp = &proc
	return
}

type GenericProcessor struct {
	procchan chan Processor

}

func (gp GenericProcessor) Process() {}

func (gp GenericProcessor) SpawnProcess(p Processor) {
	go func(){
		gp.procchan <- p
	}()
}

func (p *GenericProcessor) Init(c chan Processor) {
	p.procchan = c
}


