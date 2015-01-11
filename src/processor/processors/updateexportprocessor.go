package processors

import (
	"processor"
	"storage/types"
	"github.com/joernweissenborn/aursir4go/messages"
)

type UpdateExportProcessor struct {

	*processor.GenericProcessor

	AppId string

	UpdateExportMsg messages.UpdateExportMessage

}

func (p UpdateExportProcessor) Process() {
	                  printDebug("UPDATEEXPORT",p.UpdateExportMsg)
	Export := types.GetExportById(p.UpdateExportMsg.ExportId,p.GetAgent())
	Export.UpdateTags(p.UpdateExportMsg.Tags)
	var pjp PendingJobProcessor
	pjp.Appkey = Export.GetAppKey()
	pjp .GenericProcessor = processor.GetGenericProcessor()
	p.SpawnProcess(pjp)
	var uesp ExportedStateProcessor
	uesp.AppKey = Export.GetAppKey()
	uesp .GenericProcessor = processor.GetGenericProcessor()

	p.SpawnProcess(uesp)
}

