package boot

import "processor"

func bootCore() (processingChan chan processor.Processor){

	print("Launching Core")

	processingChan = make(chan processor.Processor)

	go processor.Process(processingChan,MAX_PROCESSORS)

	return
}
