package processors

import (
	"aursirrt/src/processor"
	"aursirrt/src/storage/types"
)

type LeaveProcessor struct {

	*processor.GenericProcessor

	AppId string

}

func (p LeaveProcessor) Process() {

	app := types.GetApp(p.AppId, p.GetAgent())

	keys := []types.AppKey{}
	for _, e:= range app.GetExports() {
		keys = append(keys,e.GetAppKey())
	}

	app.Remove()
	for _, k:= range keys {
	    if k.Exists() {
			var esp ExportedStateProcessor
			esp.AppKey = k
			esp.GenericProcessor = processor.GetGenericProcessor()
			p.SpawnProcess(esp)
		}
	}
}
