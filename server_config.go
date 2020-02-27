package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// APIToken personal access token to increase limit rate
// will be use in request's header to github API
var APIToken string

// ServerConfig structure of server.config.json file
type ServerConfig struct {
	PersonalToken string `json:"personal_token"`
}

func loadConfigFile() {
	file, err := ioutil.ReadFile("server.config.json")

	data := ServerConfig{}

	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(file), &data)

	if err != nil {
		log.Println("Error during loading config file:", err)
		return
	}

	APIToken = data.PersonalToken
}
