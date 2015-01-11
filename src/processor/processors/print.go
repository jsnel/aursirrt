package processors

import "log"

func print(msg string) {
	log.Println("PROCESSORS", msg)

}

func printDebug(msg ...interface {}) {
	if true {
		log.Println("DEBUG PROCESSORS", msg)
	}
}
