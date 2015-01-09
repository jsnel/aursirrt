package processors

import (
	"processor"
	"storage/types"
)


type DeliverJobProcessor struct {

	*processor.GenericProcessor

	Job types.Job

}

func (p DeliverJobProcessor) Process() {

	app := p.Job.GetAssignedExport().GetApp()
	app.GetConnection().Send(p.Job.GetRequest())

}
