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

	Export := types.GetExportById(p.UpdateExportMsg.ExportId,p.GetAgent())
	Export.UpdateTags(p.UpdateExportMsg.Tags)
	var pjp PendingJobProcessor
	pjp.Appkey = Export.GetAppKey()
	p.SpawnProcess(pjp)
	var uesp ExportedStateProcessor
	uesp.AppKey = Export.GetAppKey()
	p.SpawnProcess(uesp)
}

