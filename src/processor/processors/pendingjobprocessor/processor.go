package pendingjobprocessor

import (
	"processor"
	"storage/types"
	"processor/processors/deliverjobprocessor"
)


type PendingJobProcessor struct {

	*processor.GenericProcessor

	Appkey types.AppKey

}

func (p PendingJobProcessor) Process() {

	for _, imp := range p.appkey.GetImporter(){
		for _, j := range imp.GetJobs() {
			if !j.IsAssigned() {
				if imp.HasExporter() {
					exp := imp.GetExporter()[0]
					j.Assign(exp)
					var djp deliverjobprocessor.DeliverJobProcessor
					djp.Job = j
					p.SpawnProcess(djp)
				}
			}
		}
	}

}
