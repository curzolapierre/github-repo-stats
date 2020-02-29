package main

import (
	"encoding/json"
	"fmt"
)

// RepositoryGithubDto have important element from repositories list
type RepositoryGithubDto struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Owner       struct {
		OwnerName string `json:"login"`
		AvatarURL string `json:"avatar_url"`
		Type      string `json:"type"`
	}
	LanguageURL string `json:"languages_url"`
	URL         string `json:"html_url"`
}

// Repository structure of repository used and sent to the client
type Repository struct {
	Name        string
	FullName    string
	Description string
	OwnerName   string
	AvatarURL   string
	Type        string
	URL         string
	Languages   map[string]float64
}

// ContainsLanguage check if the current repository contains a searched language
func (repo *Repository) ContainsLanguage(searchedLanguage string) bool {
	for lang := range repo.Languages {
		if lang == searchedLanguage {
			return true
		}
	}
	return false
}

// FetchRepositoriesList will fetch 100 first public repositories
func FetchRepositoriesList() (*[]RepositoryGithubDto, error) {
	body, err := FetchAPI("repositories")
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch repositories", err)
	}
	repositories := &[]RepositoryGithubDto{}
	err = json.Unmarshal(body, &repositories)
	if err != nil {
		return nil, fmt.Errorf("Request returned bad JSON", err, "-", string(body))
	}

	fmt.Println("Fetched", len(*repositories), "repositories")
	return repositories, nil
}
