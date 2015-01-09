

package processors

import (
	"github.com/joernweissenborn/aursir4go"
	"processor"
	"storage/types"
	"github.com/joernweissenborn/aursir4go/messages"
)

type UpdateImportProcessor struct {

	*processor.GenericProcessor

	AppId string

	UpdateImportMsg aursir4go.AurSirUpdateImportMessage

}

func (p UpdateImportProcessor) Process() {

	Import := types.GetImportById(p.UpdateImportMsg.ImportId,p.GetAgent())
	Import.UpdateTags(p.UpdateImportMsg.Tags)
	app := Import.GetApp()
	app.Send(messages.ImportUpdatedMessage{Import.GetId(),Import.HasExporter()})
}

