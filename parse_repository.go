package main

import (
	"encoding/json"
	"fmt"
	"time"
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

type repositorySearchGithubDto struct {
	Items []RepositoryGithubDto
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

// FetchRepositoriesList will fetch 100 last public repositories
func FetchRepositoriesList(querySearch ...string) (*[]RepositoryGithubDto, error) {
	// Github API doesn't offer possibility to fetch last created repositories
	// To get them we have to set 'q' parameter, we search by date (D-1) and sort by stars arbitrary
	// to activate order desc and to be sure to load 100 last repositories
	// https://developer.github.com/v3/search/#search-repositories

	if len(querySearch) <= 1 && querySearch[0] == "" {
		yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
		querySearch = append(querySearch, "q=created:"+yesterday)
	}
	querySearch = append(querySearch, "order=desc", "sort=stars")
	body, err := FetchAPI("search/repositories", querySearch...)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch repositories", err)
	}
	repositoriesSearch := &repositorySearchGithubDto{}
	err = json.Unmarshal(body, &repositoriesSearch)
	if err != nil {
		return nil, fmt.Errorf("Request returned bad JSON", err, "-", string(body))
	}

	repositories := &repositoriesSearch.Items

	fmt.Println("Fetched", len(*repositories), "repositories")
	return repositories, nil
}
