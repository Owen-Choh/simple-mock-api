package main

import (
	"log"

	"github.com/Owen-Choh/simple-mock-api/mock-api/service"
)

var mappingsDir string = "./mappings"

func main() {
	mockServer := service.NewMockServer("/mock", mappingsDir)
	if err := mockServer.RegisterHandlers(); err != nil {
		log.Printf("error loading mappings: %s\n", err.Error())
		return
	}

	if err := mockServer.Start(":8080"); err != nil {
		log.Printf("error starting mock server: %s\n", err.Error())
		return
	}
}
