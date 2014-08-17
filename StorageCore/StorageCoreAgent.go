package storagecore

import (
	"log"
	"github.com/joernweissenborn/aursir4go"
	"github.com/joernweissenborn/aursirrt/config"
	"path"
	"os"
)

type StorageCoreAgent struct {
	write       chan StorageRequestItem
	read        chan StorageRequestItem
	storageCore StorageCore
}

func (sca *StorageCoreAgent) Launch(cfg config.RtConfig) {

	sca.read = make(chan StorageRequestItem)
	sca.write = make(chan StorageRequestItem)

	nodeid := cfg.GetConfigItem("NodeId")
	if nodeid == nil {
		log.Println("StorageCore NodeId is not found, setting")
		nodeid = generateUuid()
		cfg.SetConfigItem("NodeId",nodeid)
	}
	log.Println("StorageCore NodeId is",nodeid)
	dbpath := cfg.GetConfigItem("DatabasePath")
	if dbpath == nil {
		cwd, _ := os.Getwd()
		dbpath = path.Join(cwd,"Database")
		log.Println("StorageCore DatabasePath is not found, setting")

		cfg.SetConfigItem("DatabasePath",dbpath)
	}
	log.Println("StorageCore DatabasePath is",dbpath)
	sca.storageCore.init(nodeid.(string),dbpath.(string))

	go sca.listen()

}

func (sca StorageCoreAgent) Write(req StorageRequest) StorageReply {
	return sca.doRequest(req, sca.write)
}

func (sca StorageCoreAgent) Read(req StorageRequest) StorageReply {
	return sca.doRequest(req, sca.read)
}

func (sca StorageCoreAgent) doRequest(req StorageRequest, ch chan StorageRequestItem) StorageReply {

	replychan := make(chan StorageReply)
	ch <- StorageRequestItem{replychan, req}
	reply := <-replychan
	return reply
}

func (sca StorageCoreAgent) listen() {

	log.Println("StorageCoreAgent ready")

	for {
		select {

		case writeRequest, ok := <-sca.write:

			if ok {

				writeRequest.reply <- sca.dowrite(writeRequest)
			} else {
				writeRequest.reply <- WriteFail{}
			}

		case readRequest, ok := <-sca.read:
			if ok {
				readRequest.reply <-sca.doread(readRequest)
			}else{
		readRequest.reply <- ReadFail{}
		}

		}
	}
}

func (sca StorageCoreAgent) dowrite(req StorageRequestItem) StorageReply {
	switch request := req.request.(type) {

	case RegisterAppRequest:
		sca.storageCore.registerApp(request)
		return WriteOk{}

	case RemoveAppRequest:
		return sca.storageCore.removeApp(request)

	case AddExportRequest:
		return sca.storageCore.addExport(request)

	case AddImportRequest:
		id, exported := sca.storageCore.addImport(request)
		return ImportAdded{id, exported}

	case AddReqRequest:
		return ReqRegistered{sca.storageCore.addRequest(request)}
	case AddResRequest:
		return sca.storageCore.addResult(request)
	case ListenRequest:
		sca.storageCore.addFuncListen(request)
		return WriteOk{}
	case AddPersistentResultRequest:
		sca.storageCore.addPersistentResult(request)
		return WriteOk{}
	case UpdateExportRequest:
		return sca.storageCore.updateExport(request)
	case UpdateImportRequest:
		return sca.storageCore.updateImport(request)
	case AddCallChainRequest:
		return sca.storageCore.addCallChain(request)
	default:
		return WriteFail{}
	}

}

func (sca StorageCoreAgent) doread(req StorageRequestItem) StorageReply{
	switch request := req.request.(type) {

	case GetAppKey:
		kv := sca.storageCore.getKeyVertex(request.KeyName)
		if kv == nil {return ReadFail{}}
		k, _ :=kv.Properties.(aursir4go.AppKey)
		return k

	case GetRequest:
		rv := sca.storageCore.graph.GetVertex(request.Uuid)
		if rv == nil {return ReadFail{}}
		r, _ :=rv.Properties.(aursir4go.AurSirRequest)
		return r

	default:
		return ReadFail{}
	}
}
