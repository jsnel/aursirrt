package config

import (
	"flag"
	"fmt"
)


type connections []string

func (i *connections) String() string {
	return fmt.Sprintf("%d", *i)
}

func (i *connections) Set(value string) error {
	*i = append(*i, value )
	return nil
}

var Zconnections connections

func init(){
	flag.Var(&Zconnections, "zconnection", "e.g. 192.168.0.1:5555, if no port specified p2p will be enabled on the iface")

}

