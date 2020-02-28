package main

import (
	"encoding/json"
	"fmt"
)

func fetchLanguagesList(owner, repoName string) (map[string]float64, error) {
	body, err := FetchAPI("repos/" + owner + "/" + repoName + "/languages")
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch repositories", err)
	}

	var f interface{}
	if err := json.Unmarshal(body, &f); err != nil {
		return nil, fmt.Errorf("Error during fetching languages", err)
	}

	itemsMap := f.(map[string]interface{})

	formatedLanguages := make(map[string]float64)
	for name, size := range itemsMap {
		formatedLanguages[name] = size.(float64)
	}

	return formatedLanguages, nil
}
