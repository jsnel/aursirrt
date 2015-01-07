package pendingjobprocessor

import (
	"processor"
	"storage/types"
)


type PendingJobProcessor struct {

	*processor.GenericProcessor

	appkey types.AppKey

}

func (p PendingJobProcessor) Process() {

	for _, imp := range p.appkey.GetImporter{

	}

}
