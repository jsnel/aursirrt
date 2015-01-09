package processors

import (
	"processor"
	"storage/types"
	"github.com/joernweissenborn/aursir4go/messages"
)

type RequestProcessor struct {

	*processor.GenericProcessor

	AppId string
	Request messages.Request
}

func (p RequestProcessor) Process() {

	job := types.GetJobFromRequest(p.Request,p.GetAgent())
	job.Create()
	imp := job.GetImport()
	if job.Exists() && imp.HasExporter() {
		exp := imp.GetExporter()[0]
		job.Assign().
	}
}

