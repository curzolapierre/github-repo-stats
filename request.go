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
	req.Header.Set("Authorization", "token "+serverConfig.PersonalToken)

	return http.DefaultClient.Do(req)
}

// FetchAPI execute GET request to Github API to fetch endPoint parameter
func FetchAPI(endPoint string) ([]byte, error) {
	// TODO(pg): streaming possibilities ?
	resp, err := execRequest("GET", serverConfig.GithubAPIURL+"/"+endPoint, nil)
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
