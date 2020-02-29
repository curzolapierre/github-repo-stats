package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
	"sync"
)

// ServerConfig structure of server.config.json file
// PersonalToken APIToken personal access token to increase limit rate
// 	will be use in request's header to github API
type ServerConfig struct {
	GithubAPIURL  string `json:"github_api_url"`
	PersonalToken string `json:"personal_token"`
	WorkerNumber  int    `json:"worker_number"`
}

// will load config file from path
// If no file found no error will be returned, project can work without
func loadConfigFile(path string) (*ServerConfig, error) {
	file, err := ioutil.ReadFile(path)

	data := ServerConfig{}

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(file), &data)

	if err != nil {
		return nil, err
	}

	return &data, nil
}

// Check will check conformity of config file
func (env *ServerConfig) Check() error {
	var errorMessages []string
	errChan := make(chan error)
	wg := &sync.WaitGroup{}
	wg.Add(3)

	go checkURL(wg, errChan, env.GithubAPIURL, "github_api_url")
	go checkWorkerNumber(wg, errChan, env.WorkerNumber, "worker_number")
	go checkAPICredentials(wg, errChan, env.PersonalToken, "personal_token")

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		errorMessages = append(errorMessages, err.Error())
	}

	if len(errorMessages) > 0 {
		return errors.New("  → " + strings.Join(errorMessages, "\n  → "))
	}
	return nil
}

func checkURL(wg *sync.WaitGroup, c chan error, raw string, name string) {
	defer wg.Done()
	_, err := url.Parse(raw)
	if err != nil || raw == "" {
		c <- fmt.Errorf("%v is not a valid URL: %v", name, err)
	}
}

func checkWorkerNumber(wg *sync.WaitGroup, c chan error, raw int, name string) {
	defer wg.Done()
	if raw <= 0 {
		c <- fmt.Errorf("%v is not valid, need to be greater then 0", name)
	}
}

func checkAPICredentials(wg *sync.WaitGroup, c chan error, raw string, name string) {
	defer wg.Done()
	if raw == "" {
		c <- fmt.Errorf("%v is empty, you can increase the rate limit when this field is filled in. See README.md", name)
	}
}
