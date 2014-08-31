package webserver

import (
	"github.com/joernweissenborn/aursirrt/config"
	"os"
	"path"
	"log"
	"net/http"
)

func Launch(cfg config.RtConfig) {

	log.Println("WEBSERVER","Launching")


	webfolder:= cfg.GetConfigItem("Webfolder")
	if webfolder == nil{
		cwd, _ := os.Getwd()
		webfolder = path.Join(cwd,"www")
		log.Println("WEBSERVER","Webfolder is not found, setting",webfolder)

		cfg.SetConfigItem("Webfolder",webfolder)

	}

	log.Fatal("WEBSERVER",http.ListenAndServe("0.0.0.0:8080", http.FileServer(http.Dir(webfolder.(string)))))

}
