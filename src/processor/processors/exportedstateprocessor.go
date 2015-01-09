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
		conn := imp.GetApp().GetConnection()
		conn.Send(messages.ImportUpdatedMessage{imp.GetId(),imp.HasExporter()})
	}
}

