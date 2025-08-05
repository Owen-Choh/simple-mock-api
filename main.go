package main

import (
	"fmt"
	"github.com/Owen-Choh/simple-mock-api/mock-api/types"
	"github.com/Owen-Choh/simple-mock-api/mock-api/utils"
	"os"
)

func main() {
	var config types.Response
	configFile, err := os.Open("config.json")
	if err != nil {
		fmt.Printf("opening config file: %s\n", err.Error())
	}

	if err := utils.ParseJsonFile(configFile, &config); err != nil {
		fmt.Printf("parsing config file: %s\n", err.Error())
		return
	}

	fmt.Printf("%d %s", config.Status, config.Message)
	return
}
