package datastorage

import "github.com/joernweissenborn/aursirrt/config"

import (
	"github.com/joernweissenborn/aursir4go"
	"os"
	"log"
	"path"
)




func Open(cfg config.RtConfig, requestChan chan interface{}) {

	storagepath:= cfg.GetConfigItem("StoragePath")
	if storagepath == nil{
		cwd, _ := os.Getwd()
		storagepath = path.Join(cwd,"Database")
		log.Println("StorageCore DatabasePath is not found, setting")

		cfg.SetConfigItem("DatabasePath",storagepath)

	}
	CreateFolderIfNotExist(storagepath.(string))
	for r := range requestChan {
		switch req := r.(type) {
			case CommitData:
				commitData(req)
		}
	}

}

func CreateFolderIfNotExist(datapath string) error {
	dir, err := os.Open(datapath)
	if err != nil {
			err =os.MkdirAll(datapath,os.ModeDir)
			if err != nil {
				return err
			}
	}
	dir.Close()
	return nil
}

type CommitData struct {
	Request aursir4go.AurSirRequest
	Result aursir4go.AurSirResult
}

func commitData(data CommitData){

}
