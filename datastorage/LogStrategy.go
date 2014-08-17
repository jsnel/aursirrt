package datastorage

import (
	"log"
	"os"
	"fmt"
	"path"
)

type logPersistor struct {}

func (logPersistor) PersistData(filepath string,data CommitData) error{
	f, err := os.OpenFile(path.Join(filepath,data.Result.FunctionName+".tsv"), os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Println("LogPersistor","error opening file: %v", err)
		return err
	}
	defer f.Close()

	entry := fmt.Sprintf("%s\t%s\t%s\t%s\n",data.Result.Timestamp,data.Result.Uuid,data.Request.Request,data.Result.Result)

	_,err = f.Write([]byte(entry))

	return err
}
