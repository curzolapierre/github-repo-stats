package main

import (
	"log"
	"net/http"
)

var serverConfig *ServerConfig

func main() {
	localServerConfig, err := loadConfigFile("./server.config.json")
	if err != nil {
		log.Fatalln("Error during loading config file:", err)
	}
	cred, err := loadConfigFile("./credentials.json")
	if err != nil {
		log.Print("Error during loading credentials file:", err)
	} else {
		localServerConfig.PersonalToken = cred.PersonalToken
		err = localServerConfig.Check()
		if err != nil {
			log.Fatalf("Errors in config files:\n%v\n", err)
		}
	}
	serverConfig = localServerConfig

	mux := http.NewServeMux()
	mux.HandleFunc("/", makeHandler(indexHandler))
	mux.HandleFunc("/index", makeHandler(indexHandler))
	mux.HandleFunc("/search/", makeHandler(searchHandler))

	log.Println("Starting server on :8080...")
	err = http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
