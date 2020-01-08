package main

import (
	"github.com/DerekKeeler/server-finder-demo/server"
)

func main() {
	server.Start(":8080", "Test")
}
