package main

import (
	"flag"
	"go-file-server/routes"
	"go-file-server/utils"
	"log"
	"net/http"
)

func main() {
	// Get parameters
	flag.Parse()

	// Initialize server
	http.HandleFunc("/", routes.FileServer)
	log.Printf("Serving %s on HTTP port: %s\n", *utils.Directory, *utils.Port)
	log.Fatal(http.ListenAndServe(":"+*utils.Port, nil))
}
