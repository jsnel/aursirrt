package StorageCore

import "log"

type StorageCoreAgent struct {
	write       chan StorageRequestItem
	read        chan StorageRequestItem
	storageCore StorageCore
}

func (sca *StorageCoreAgent) Launch() {

	sca.read = make(chan StorageRequestItem)
	sca.write = make(chan StorageRequestItem)

	sca.storageCore.init()

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
				sca.dowrite(readRequest)
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
		return ResRegistered{sca.storageCore.addResult(request)}
	case ListenRequest:
		sca.storageCore.addFuncListen(request)
		return WriteOk{}
	default:
		return WriteFail{}
	}

}

func (sca StorageCoreAgent) doread(req StorageRequestItem) {

}
