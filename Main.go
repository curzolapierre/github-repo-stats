package main

import (
	"fmt"
	"log"
)

var serverConfig *ServerConfig

func main() {
	localServerConfig, err := loadConfigFile("./server.config.json")
	if err != nil {
		log.Fatalln("Error during loading config file:", err)
	} else {
		err = localServerConfig.Check()
		if err != nil {
			log.Fatalf("Errors in config file:\n%v\n", err)
		}
		fmt.Println("github URL: ", localServerConfig.GithubAPIURL)
	}
	serverConfig = localServerConfig
	// quitCh := make(chan struct{})

	_, err = getAggregatedRepo()
	if err != nil {
		log.Fatalln(err)
	}

	// mux := http.NewServeMux()

	// log.Println("Starting server on :8080...")
}
