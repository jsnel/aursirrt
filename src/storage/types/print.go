package types

import "log"

func print(msg string) {
	log.Println("STORAGECORE", msg)

}

func printDebug(msg string) {
	if true {
		log.Println("DEBUG STORAGECORE", msg)
	}
}
