package processors

import (
	"aursirrt/src/processor"
	"aursirrt/src/storage/types"
	"github.com/joernweissenborn/aursir4go/messages"
)

type StartListenProcessor struct {

	*processor.GenericProcessor

	AppId string

	StartListenMsg messages.ListenMessage

}

func (p StartListenProcessor) Process() {

	Import := types.GetImportById(p.StartListenMsg.ImportId,p.GetAgent())

	Import.StartListenToFunction(p.StartListenMsg.FunctionName)




}

