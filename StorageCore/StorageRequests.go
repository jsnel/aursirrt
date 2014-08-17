package storagecore

import "github.com/joernweissenborn/aursir4go"

//Base interface for all storage requests
type StorageRequestItem struct {
	reply   chan StorageReply
	request StorageRequest
}

type StorageRequest interface{}

type RegisterAppRequest struct {
	Id      string
	AppName string
}

type RemoveAppRequest struct {
	Id string
}

type AddExportRequest struct {
	Id     string
	AppKey aursir4go.AppKey
	Tags   []string
}

type UpdateExportRequest struct {
	Req aursir4go.AurSirUpdateExportMessage
}

type UpdateImportRequest struct {
	Req aursir4go.AurSirUpdateImportMessage
}


type AddImportRequest struct {
	Id     string
	AppKey aursir4go.AppKey
	Tags   []string
}

type AddReqRequest struct {
	AppId string
	Req aursir4go.AurSirRequest
}

type AddResRequest struct {
	AppId string
	Req aursir4go.AurSirResult
}

type AddPersistentResultRequest struct {
	Req aursir4go.AurSirResult
}

type AddCallChainRequest struct {
	AppId string
	Req aursir4go.AurSirCallChain
}

type ListenRequest struct {
	AppId string
	FuncName string
	ImportId string
}

type GetAppKey struct {
	KeyName string
}

type GetRequest struct {
	Uuid string
}

//StorageReply is the base interfaces for all replies to storageRequests
type StorageReply interface{}

type WriteOk struct {
}


type WriteFail struct {
}
type ReadFail struct {
}

type ExportAdded struct {
	ExportId         string
	ConnectedImports map[string]string
	DisconnectedImports map[string]string
	PendingJobs []aursir4go.AurSirRequest
}

type ImportAdded struct {
	ImportId string
	Exported bool
}

type AppRemoved struct {
	DisconnectedImports map[string]string
}

type ReqRegistered struct {
	Exporter []string
}
type ResRegistered struct {
	Importer []string
	IsChainCall bool
	ChainCall aursir4go.ChainCall
	ChainCallImportId string
}
