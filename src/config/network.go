package config

import "flag"


var Zmqport = flag.Int("zport", 5555, "Set custom port for zeromq backend")
var Zmqip = flag.String("zip", "localhost", "Set custom ip for zeromq backend")
