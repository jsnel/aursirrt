package StorageCore

import "github.com/joernweissenborn/AurSir4Go"

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
	AppKey AurSir4Go.AppKey
	Tags   []string
}
type AddImportRequest struct {
	Id     string
	AppKey AurSir4Go.AppKey
	Tags   []string
}

type AddReqRequest struct {
	AppId string
	Req AurSir4Go.AurSirRequest
}

type AddResRequest struct {
	AppId string
	Req AurSir4Go.AurSirResult
}

type ListenRequest struct {
	AppId string
	FuncName string
	ImportId string
}

//StorageReply is the base interfaces for all replies to storageRequests
type StorageReply interface{}

type WriteOk struct {
}

type WriteFail struct {
}

type ExportAdded struct {
	ExportId         string
	ConnectedImports map[string]string
	PendingJobs []AurSir4Go.AurSirRequest
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
}
