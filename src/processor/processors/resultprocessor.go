package processors

import (
	"processor"
	"storage/types"
	"github.com/joernweissenborn/aursir4go/messages"
	"github.com/joernweissenborn/aursir4go/calltypes"
	"log"
)

type ResultProcessor struct {

	*processor.GenericProcessor

	AppId string
	Result messages.Result
}

func (p ResultProcessor) Process() {

	job := types.GetJobFromResult(p.Result,p.GetAgent())
	               	log.Println(p.Result)
	switch p.Result.CallType {
	case calltypes.ONE2ONE, calltypes.ONE2MANY:
		if job.Exists() {
			imp := job.GetImport()
			var smp SendMessageProcessor
			smp.App = imp.GetApp()
			smp.Msg = p.Result
			smp.GenericProcessor = processor.GetGenericProcessor()
			p.SpawnProcess(smp)
			job.Remove()
		}
	case calltypes.MANY2MANY, calltypes.MANY2ONE:
		exp := types.GetExportById(p.Result.ExportId,p.GetAgent())
	for _,imp := range exp.GetAppKey().GetImporter()   {

		if exp.HasTags(imp.GetTagNames()){
			var smp SendMessageProcessor
			smp.App = imp.GetApp()
			smp.Msg = p.Result
			smp.GenericProcessor = processor.GetGenericProcessor()
			p.SpawnProcess(smp)
		}
	}

	}
}

