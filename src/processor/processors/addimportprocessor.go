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

	Import := types.GetImport(p.AppId,p.AddImportMsg.AppKey, p.AddImportMsg.Tags ,p.AddImportMsg.ImportId,p.GetAgent())
	Import.Add()
	var smp SendMessageProcessor
	smp.App = Import.GetApp()
	smp.Msg = messages.ImportAddedMessage{Import.GetId(),Import.HasExporter()}
	smp.GenericProcessor = processor.GetGenericProcessor()
	p.SpawnProcess(smp)



}

