package addexportprocessor

import (
	"github.com/joernweissenborn/aursir4go"
	"processor"
	"storage/types"
)

type AddExportProcessor struct {

	*processor.GenericProcessor

	AppId string

	AddExportMsg aursir4go.AurSirAddExportMessage

}

func (p AddExportProcessor) Process() {

	export := types.GetExport(p.AppId,p.AddExportMsg.AppKey, p.AddExportMsg.Tags,p.GetAgent())
	export.Add()

}

