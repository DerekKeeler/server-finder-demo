package main

import (
	"log"

	"github.com/DerekKeeler/server-finder-demo/client"
)

func main() {
	responses, err := client.Scan()

	if err != nil {
		log.Println(err)
	}

	log.Printf("%+v", responses)
}
