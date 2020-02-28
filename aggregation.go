package main

import "fmt"

// Pipeline:
// 3 stages
// 1st: list repo to chan
// 2nd: fetch languages

// Manager will fetch public repositories, their languages then call worker to aggregate those data
func Manager() {
	repoList, err := fetchRepositoriesList()
	if err != nil {
		return
	}

	// Set up the pipeline.
	for repo := range fetchAndTransferRepoToTheChan(listRepoToTheChan(*repoList)) {
		fmt.Println(repo)
	}
}

func listRepoToTheChan(repoList []RepositoryGithubDto) <-chan RepositoryGithubDto {
	out := make(chan RepositoryGithubDto)
	go func() {
		for _, repo := range repoList {
			out <- repo
		}
		close(out)
	}()
	return out
}

func fetchAndTransferRepoToTheChan(in <-chan RepositoryGithubDto) <-chan Repositories {
	out := make(chan Repositories)
	go func() {
		for repoDto := range in {
			languages, err := fetchLanguagesList(repoDto.Owner.OwnerName, repoDto.Name)
			if err != nil {
				return
			}

			repo := Repositories{
				Name:        repoDto.Name,
				FullName:    repoDto.FullName,
				Description: repoDto.Description,
				OwnerName:   repoDto.Owner.OwnerName,
				AvatarURL:   repoDto.Owner.AvatarURL,
				Type:        repoDto.Owner.Type,
				URL:         repoDto.URL,
				Languages:   languages,
			}

			out <- repo
		}
		close(out)
	}()
	return out
}
