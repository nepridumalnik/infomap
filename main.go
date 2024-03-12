package main

import (
	"infomap/app"
	"log"
)

func main() {
	app, err := app.CreateApp(":8080")

	if err != nil {
		log.Fatal(err)
	}

	err = app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
