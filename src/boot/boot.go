package boot

import (
	"log"
)

const (
	MAX_PROCESSORS = 4
)

func Boot(){

	print("AurSir RT starting")

	BootCore()

}

func print(msg string){

	log.Println("BOOT",msg)

}
