package tags

import (
	"storage"
	"storage/types/appkey"
)

type Tag struct {
	agent storage.StorageAgent
	key appkey.AppKey
	name string
}

func Get(key appkey.AppKey, name string, agent storage.StorageAgent)Tag{
	return Tag{agent,key,name}
}

func Add(){

}
