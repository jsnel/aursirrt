package addexportprocessor

import (
	"github.com/joernweissenborn/aursir4go"
	"processor"
	"storage/types"
	"github.com/joernweissenborn/aursir4go/messages"
	"processor/processors/pendingjobprocessor"
)

type AddExportProcessor struct {

	*processor.GenericProcessor

	AppId string

	AddExportMsg aursir4go.AurSirAddExportMessage

}

func (p AddExportProcessor) Process() {
	app := types.GetApp(p.AppId, p.GetAgent())

	if !app.Exist(){
		return
	}

	export := types.GetExport(p.AppId,p.AddExportMsg.AppKey, p.AddExportMsg.Tags,p.GetAgent())
	export.Add()

	app.GetConnection().Send(messages.ExportAddedMessage{export.GetId()})
	var pjp pendingjobprocessor.PendingJobProcessor
	pjp.Appkey = export.GetAppKey()
	p.SpawnProcess(pjp)
}

