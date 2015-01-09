package processors

import (
	"processor"
	"storage/types"
	"github.com/joernweissenborn/aursir4go/messages"
)

type AddExportProcessor struct {

	*processor.GenericProcessor

	AppId string

	AddExportMsg messages.AddExportMessage

}

func (p AddExportProcessor) Process() {

	Export := types.GetExport(p.AppId,p.AddExportMsg.AppKey,p.AddExportMsg.Tags,p.GetAgent())
	Export.Add()
	app := Export.GetApp()
	app.Send(messages.ExportAddedMessage{Export.GetId()})
	var pjp PendingJobProcessor
	pjp.Appkey = Export.GetAppKey()
	p.SpawnProcess(pjp)
	var uesp ExportedStateProcessor
	uesp.AppKey = Export.GetAppKey()
	p.SpawnProcess(uesp)
}

