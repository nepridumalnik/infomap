package main

import (
	server "infomap/src"
	"log"
)

func main() {
	server, err := server.CreateServer()

	if err != nil {
		log.Fatal(err)
	}

	err = server.ListenAndServe(":8080")

	if err != nil {
		log.Fatal(err)
	}
}
