package main

func manageLanguages(repoList *[]RepositoryGithubDto) {
	for _, repo := range *repoList {
		languages, err := fetchLanguagesList(repo.Owner.OwnerName, repo.Name)
		if err != nil {
			return
		}
		repositoryObj := &Repositories{
			Name:        repo.Name,
			FullName:    repo.FullName,
			Description: repo.Description,
			OwnerName:   repo.Owner.OwnerName,
			AvatarURL:   repo.Owner.AvatarURL,
			Type:        repo.Owner.Type,
			URL:         repo.URL,
			Languages:   languages,
		}
		aggregateRepo(repositoryObj)
	}
}
