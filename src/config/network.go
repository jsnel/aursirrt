package config

import "flag"


var Zmqport = flag.Int("zport", 5555, "Set custom port for zeromq backend")
var Zmqip = flag.String("zip", "localhost", "Set custom ip for zeromq backend")
var P2p = flag.Bool("p2p", false, "Set to enables runtime p2p connection.")
