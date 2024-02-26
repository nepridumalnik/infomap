package main

import (
	server "infomap/src"
)

func main() {
	server := server.CreateServer()
	server.ListenAndServe(":8080")
}
