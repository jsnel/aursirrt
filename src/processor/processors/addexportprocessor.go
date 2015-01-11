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
	             printDebug("AddExport")
	Export := types.GetExport(p.AppId,p.AddExportMsg.AppKey,p.AddExportMsg.Tags,p.GetAgent())
	Export.Add()
	app := Export.GetApp()
	var smp SendMessageProcessor
	smp.App = app
	smp.Msg = messages.ExportAddedMessage{Export.GetId()}
	smp.GenericProcessor = processor.GetGenericProcessor()
	p.SpawnProcess(smp)
	var pjp PendingJobProcessor
	pjp.Appkey = Export.GetAppKey()
	pjp.GenericProcessor = processor.GetGenericProcessor()
	p.SpawnProcess(pjp)
	var uesp ExportedStateProcessor
	uesp.AppKey = Export.GetAppKey()
	uesp.GenericProcessor = processor.GetGenericProcessor()
	p.SpawnProcess(uesp)
}

