package main

import (
	"fmt"
	"log"
	"net/http"
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

	mux := http.NewServeMux()
	mux.HandleFunc("/", makeHandler(indexHandler))
	mux.HandleFunc("/index", makeHandler(indexHandler))
	mux.HandleFunc("/search/", makeHandler(searchHandler))
	mux.HandleFunc("/query/", makeHandler(queryHandler))

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Starting server on :8080...")
	err = http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
