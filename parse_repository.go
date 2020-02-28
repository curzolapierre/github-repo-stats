package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func fetchRepositoriesList(w http.ResponseWriter, r *http.Request) (*[]RepositoryGithubDto, error) {
	body, err := FetchAPI("repositories")
	if err != nil {
		fmt.Errorf("Failed to fetch repositories", err)
		return nil, err
	}
	repositories := &[]RepositoryGithubDto{}
	err = json.Unmarshal(body, &repositories)
	if err != nil {
		fmt.Errorf("Request returned bad JSON", err, "-", string(body))
		return nil, err
	}
	owner := (*repositories)[1].Owner.OwnerName
	name := (*repositories)[1].Name
	fmt.Println("Repository:")
	fmt.Println("  → Name", (*repositories)[0].Name)
	fmt.Println("  → FullName", (*repositories)[0].FullName)
	fmt.Println("  → Description", (*repositories)[0].Description)
	fmt.Println("  → AvatarURL", (*repositories)[0].Owner.AvatarURL)
	fmt.Println("  → LanguageURL", (*repositories)[0].LanguageURL)
	fmt.Println("  → OwnerName", (*repositories)[0].Owner.OwnerName)
	fmt.Println("  → Type", (*repositories)[0].Owner.Type)
	fmt.Println("  → URL", (*repositories)[0].URL)

	languages, err := fetchLanguagesList(owner, name)
	if err != nil {
		fmt.Errorf("Error during fetching languages", err)
		return nil, err
	}
	for name, size := range languages {
		fmt.Println("languages:", name, size)
	}
	return repositories, nil
}
