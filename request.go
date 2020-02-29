package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func execRequest(method string, url string, payload interface{}) (*http.Response, error) {
	payloadJSON, err := json.Marshal(&payload)
	if err != nil {
		return nil, fmt.Errorf("Fail to encode to JSON: %v", err)
	}
	payloadBuffer := bytes.NewBuffer(payloadJSON)

	req, err := http.NewRequest(method, url, payloadBuffer)
	if err != nil {
		return nil, err
	}

	if req.Method != "DELETE" {
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")
	}
	if serverConfig.PersonalToken != "" {
		req.Header.Set("Authorization", "token "+serverConfig.PersonalToken)
	}

	return http.DefaultClient.Do(req)
}

// FetchAPI execute GET request to Github API to fetch endPoint parameter
// params each element has to be already formated like key=value
func FetchAPI(endPoint string, params ...string) ([]byte, error) {

	queryParam := "?"
	for _, param := range params {
		queryParam += param
	}

	url := serverConfig.GithubAPIURL + "/" + endPoint
	if len(params) > 0 {
		url += queryParam
	}

	// TODO(pg): streaming possibilities ?
	resp, err := execRequest("GET", url, params)
	if err != nil {
		return nil, fmt.Errorf("Fail to get resource", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 && resp.StatusCode != 201 && resp.StatusCode != 202 {
		return nil, fmt.Errorf("Request returned bad status", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read body from request", err)
	}
	remaining := resp.Header.Get("X-RateLimit-Remaining")
	fmt.Println("requests remaining:", remaining)
	return body, nil
}
