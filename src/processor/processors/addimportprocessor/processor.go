package addimportprocessor

import (
	"github.com/joernweissenborn/aursir4go"
	"processor"
	"storage/types"
)

type AddImportProcessor struct {

	*processor.GenericProcessor

	AppId string

	AddImportMsg aursir4go.AurSirAddImportMessage

}

func (p AddImportProcessor) Process() {

	Import := types.GetImport(p.AppId,p.AddImportMsg.AppKey, p.AddImportMsg.Tags,p.GetAgent())
	Import.Add()

}

