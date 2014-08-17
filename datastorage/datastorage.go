package datastorage

import "github.com/joernweissenborn/aursirrt/config"

import (
	"github.com/joernweissenborn/aursir4go"
	"os"
	"log"
	"path"
	"fmt"
)




func Open(cfg config.RtConfig, requestChan chan interface{}) {

	storagepath:= cfg.GetConfigItem("DataStorePath")
	if storagepath == nil{
		cwd, _ := os.Getwd()
		storagepath = path.Join(cwd,"DataStorePath")
		log.Println("DATASTORAGE","DatastorPath is not found, setting")

		cfg.SetConfigItem("DataStorePath",storagepath)

	}
	createFolderIfNotExist(storagepath.(string))
	for r := range requestChan {
		switch req := r.(type) {
			case CommitRequest:
				req.Answer <- commitData(storagepath.(string),req.Data)
		}
	}

}

func createFolderIfNotExist(datapath string) error {
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

type CommitRequest struct {
	Answer chan string
	Data CommitData
}

type CommitData struct {
	Request *aursir4go.AurSirRequest
	Result *aursir4go.AurSirResult
}

func commitData(basepath string,data CommitData) (path string){
	path = createDataPath(basepath,data.Result)
	err := createFolderIfNotExist(path)
	if err != nil {
		log.Println("DATASTORAGE","Error creating path",err)
	}
	persistor := getPersistor(data.Result.PersistenceStrategy)
	persistor.PersistData(path,data)
	return
}

func createDataPath(basepath string, result *aursir4go.AurSirResult) string{
	datestring := fmt.Sprintf("%d",result.Timestamp.Year())+fmt.Sprintf("%02d",result.Timestamp.Month())+
			fmt.Sprintf("%02d",result.Timestamp.Day())
	keyandtagsstring := result.AppKeyName
	for _,tag := range result.Tags {
		keyandtagsstring = keyandtagsstring+"@"+tag
	}
	return path.Join(basepath,datestring,keyandtagsstring)
}
