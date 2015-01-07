package types

import (
	"storage"
	"fmt"
	"github.com/joernweissenborn/aursir4go/messages"
)

type jobproperties struct {
	request *messages.Request
	result *messages.Result
}
type Job struct {
	agent storage.StorageAgent
	request *messages.Request
	result *messages.Result
}

func GetJobFromRequest(request messages.Request, agent storage.StorageAgent)(Job Job){
	Job.agent = agent
	Job.request = &request
	return
}

func GetJobFromResult(result messages.Result, agent storage.StorageAgent)(Job Job){
	Job.agent = agent
	Job.result = &result
	return
}
func GetJobById(id string, agent storage.StorageAgent)(Job Job){
	Job.agent = agent
	Job.result = &result
	c := make(chan bool)
	defer close(c)
	j.agent.Read(func (sc *storage.StorageCore){
		kv := sc.GetVertex(j.request.Uuid)
		if kv == nil {
			c <- Job
			return
		}
		c<-true
		return
	})
	return
}

func (j *Job) GetId() string{

	return j.request.Uuid

}

func (j *Job) Create(){
	imp := GetImportById(j.request.ImportId,j.agent)
	if imp.Exists() {
		printDebug(fmt.Sprint("Creating Job"))

		c := make(chan string)
		defer close(c)

		j.agent.Write(func (sc *storage.StorageCore){

			jv := sc.CreateVertex(j.request.Uuid, jobproperties{j.request,j.result})
			iv := sc.GetVertex(j.request.ImportId)
			sc.CreateEdge(storage.GenerateUuid(), storage.AWAITING_JOB_EDGE, jv,iv, nil)

			c <- ""
		})

	}
}


func (j Job) Exists() bool {
	c := make(chan bool)
	defer close(c)
	j.agent.Read(func (sc *storage.StorageCore){
		kv := sc.GetVertex(j.request.Uuid)
		if kv == nil {
			c <-false
			return
		}
		c<-true
		return
	})
	return <- c
}
