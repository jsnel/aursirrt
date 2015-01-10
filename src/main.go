package main

import (
	"boot"
	"cmdlineinterface"
)

func main() {
	boot.Boot()
	cli := cmdlineinterface.CmdLineInterface{}
	cli.Run()
}

