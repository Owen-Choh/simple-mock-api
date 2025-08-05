package main

import (
	"log"
	"os"

	"github.com/Owen-Choh/simple-mock-api/mock-api/service"
	"github.com/Owen-Choh/simple-mock-api/mock-api/types"
	"github.com/Owen-Choh/simple-mock-api/mock-api/utils"
)

func main() {
	var mappings []types.Mapping
	mappingsFile, err := os.Open("config.json")
	if err != nil {
		log.Printf("opening config file: %s\n", err.Error())
		return
	}
	defer mappingsFile.Close()

	if err := utils.ParseJsonFile(mappingsFile, &mappings); err != nil {
		log.Printf("parsing config file: %s\n", err.Error())
		return
	}

	for _, mapping := range mappings {
		log.Printf("path: %s\nmethod: %s\nstatus code: %d\nbody: %s\n",
			mapping.Path, mapping.Method, mapping.Response.StatusCode, mapping.Response.Body)
	}

	mockServer := service.NewMockServer()
	if err := mockServer.LoadMappings("/mock", mappings); err != nil {
		log.Printf("error loading mappings: %s\n", err.Error())
		return
	}

	if err := mockServer.Start(":8080"); err != nil {
		log.Printf("error starting mock server: %s\n", err.Error())
		return
	}
}
