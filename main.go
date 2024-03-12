package main

import (
	server "infomap/app"
	"log"
)

func main() {
	app, err := server.CreateApp(":8080")

	if err != nil {
		log.Fatal(err)
	}

	err = app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
