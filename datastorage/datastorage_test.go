package datastorage

import (
	"testing"
	"github.com/joernweissenborn/aursir4go"
	"time"
)

func TestCommit(t *testing.T){
	var req aursir4go.AurSirRequest
	var res aursir4go.AurSirResult
	req.Uuid = time.Now().String()
	commitData(CommitData{req,res})

}
