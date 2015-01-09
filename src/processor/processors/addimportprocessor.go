package processors

import (
	"processor"
	"storage/types"
	"github.com/joernweissenborn/aursir4go/messages"
)

type AddImportProcessor struct {

	*processor.GenericProcessor

	AppId string

	AddImportMsg messages.AddImportMessage

}

func (p AddImportProcessor) Process() {

	Import := types.GetImport(p.AppId,p.AddImportMsg.AppKey, p.AddImportMsg.Tags,p.GetAgent())
	Import.Add()
	Import.GetApp().Send(messages.ImportAddedMessage{Import.GetId(),Import.HasExporter()})


}

