package main

import (
	"fmt"
)

var serverConfig *ServerConfig

func main() {
	localServerConfig, err := loadConfigFile("./server.config.json")
	if err != nil {
		fmt.Println("Warning during loading config file:", err)
		fmt.Println("No Api token found, rate limit is fixed to 60 request per hour")
	} else {
		err = localServerConfig.Check()
		if err != nil {
			fmt.Printf("Errors in config file:\n%v\n", err)
			return
		}
		fmt.Println("Api token: ", localServerConfig.PersonalToken)
		fmt.Println("github URL: ", localServerConfig.GithubAPIURL)
	}
	serverConfig = localServerConfig
	quitCh := make(chan struct{})

	Manager()

	quitCh <- struct{}{}
}
