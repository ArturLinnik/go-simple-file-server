package utils

import "flag"

var Port = flag.String("p", "8100", "port to serve on")
var Directory = flag.String("d", ".", "the directory to host")
