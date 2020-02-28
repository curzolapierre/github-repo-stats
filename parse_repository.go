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

// Repositories structure of repository used and sent to the client
type Repositories struct {
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
func (repo *Repositories) ContainsLanguage(searchedLanguage string) bool {
	for lang := range repo.Languages {
		if lang == searchedLanguage {
			return true
		}
	}
	return false
}

func fetchRepositoriesList() (*[]RepositoryGithubDto, error) {
	body, err := FetchAPI("repositories")
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch repositories", err)
	}
	repositories := &[]RepositoryGithubDto{}
	err = json.Unmarshal(body, &repositories)
	if err != nil {
		return nil, fmt.Errorf("Request returned bad JSON", err, "-", string(body))
	}

	return repositories, nil
}
