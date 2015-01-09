package config

import (
	"os"
	"encoding/json"
	"path"
)

func NewRtCfg() (cfg RtConfig,err error){
	cwd, err:= os.Getwd()
	if err != nil {
		return
	}

	cfgpath := path.Join(cwd, "config.json")
	err = cfg.Open(cfgpath)
	return
}

type RtConfig struct {
	file string
}

func (cfg *RtConfig) Open(path string) (err error){
	file, err := os.OpenFile(path,os.O_RDWR,0666)
	if err != nil {
		file,err = os.Create(path)
	}
	cfg.file = file.Name()
	file.Close()

	return
}

func (cfg *RtConfig) GetConfigItem(name string) (item interface {}){

	item = cfg.getAllConfigItems()[name]

	return
}

func (cfg *RtConfig) SetConfigItem(name string, value interface {}){
	items := cfg.getAllConfigItems()
	items[name] = value

	file,_ := os.Create(cfg.file)
	j, _ := json.MarshalIndent(items, "", "  ")
	///enc := json.NewEncoder(file)
	//enc.Encode(items)
	file.Write(j)
	file.Close()
}
func (cfg *RtConfig) getAllConfigItems() map[string] interface {}{
	file, _ := os.Open(cfg.file)
	dec := json.NewDecoder(file)
	var f interface {}
	err := dec.Decode(&f)
	file.Close()
	m := map[string]interface {}{}
	if err ==nil {
		m = f.(map[string]interface{})
	}
	return m
}
