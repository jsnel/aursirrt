package datastorage

import (
	"testing"
	"github.com/joernweissenborn/aursir4go"
	"time"
	"os"
)

func TestCommit(t *testing.T){
	var req aursir4go.AurSirRequest
	var res aursir4go.AurSirResult
	res.Uuid = time.Now().String()
	res.Timestamp = time.Now()
	req.Request = []byte("{request}")
	res.Result = []byte("{result}")
	//commitData(CommitData{req,res})

}

func TestLogging(t *testing.T){
	var l logPersistor
	var req aursir4go.AurSirRequest
	var res aursir4go.AurSirResult
	res.Uuid = "1234"
	res.Timestamp = time.Now()
	req.Request = []byte("{request}")
	res.Result = []byte("{result}")
	wd, _ := os.Getwd()
	l.PersistData(wd,CommitData{&req,&res})


}
