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
}

// will load config file into path parameter
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
	wg.Add(1)

	go checkURL(wg, errChan, env.GithubAPIURL, "github_api_url")

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
	if err != nil {
		c <- fmt.Errorf("%v is not a valid URL: %v", name, err)
	}
}
