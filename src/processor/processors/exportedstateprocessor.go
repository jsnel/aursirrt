package processors

import (
	"processor"
	"storage/types"

	"github.com/joernweissenborn/aursir4go/messages"
)

type ExportedStateProcessor struct {

	*processor.GenericProcessor

	AppKey types.AppKey


}

func (p ExportedStateProcessor) Process() {
	for _, imp := range p.AppKey.GetImporter() {
		var smp SendMessageProcessor
		smp.App = imp.GetApp()
		smp.Msg =messages.ImportUpdatedMessage{imp.GetId(),imp.HasExporter()}
		smp.GenericProcessor = processor.GetGenericProcessor()
		p.SpawnProcess(smp)
	}
}

