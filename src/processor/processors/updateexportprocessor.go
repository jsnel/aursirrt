package processors

import (
	"github.com/joernweissenborn/aursir4go"
	"processor"
	"storage/types"
	"github.com/joernweissenborn/aursir4go/messages"
)

type UpdateExportProcessor struct {

	*processor.GenericProcessor

	AppId string

	UpdateExportMsg aursir4go.AurSirUpdateExportMessage

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

