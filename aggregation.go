package main

import (
	"fmt"
	"sync"
	"time"
)

type result struct {
	repository Repository
	err        error
}

// Pipeline:
// 2 stages
// 1st: list repo to chan
// 2nd: fetch languages

// getAggregatedRepo will fetch public repositories, their languages then call worker to aggregate those data
func getAggregatedRepo() (map[string]Repository, error) {
	start := time.Now()

	repoDtoList, err := FetchRepositoriesList()
	if err != nil {
		return nil, err
	}

	done := make(chan struct{})
	defer close(done)

	repoCh := loopThroughRepo(done, *repoDtoList)
	c := make(chan result)
	var wg sync.WaitGroup
	const numWorker = 10
	wg.Add(numWorker)
	for i := 0; i < numWorker; i++ {
		go func() {
			worker(done, repoCh, c)
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(c)
	}()
	// End of pipeline.

	repoList := make(map[string]Repository)
	for r := range c {
		if r.err != nil {
			fmt.Println("Error detect during fetching languages from repository", err)
			continue
		}
		fmt.Println(r.repository)
		repoList[r.repository.FullName] = r.repository
	}

	elapsed := time.Since(start)
	fmt.Printf("Fetch repositories with their languages list took %s", elapsed)
	return repoList, nil
}

func worker(done <-chan struct{}, repositoriesDto <-chan RepositoryGithubDto, c chan<- result) {
	for repoDto := range repositoriesDto {
		languages, err := FetchLanguagesList(repoDto.Owner.OwnerName, repoDto.Name)
		if err != nil {
			return
		}
		repo := Repository{
			Name:        repoDto.Name,
			FullName:    repoDto.FullName,
			Description: repoDto.Description,
			OwnerName:   repoDto.Owner.OwnerName,
			AvatarURL:   repoDto.Owner.AvatarURL,
			Type:        repoDto.Owner.Type,
			URL:         repoDto.URL,
			Languages:   languages,
		}

		select {
		case c <- result{repository: repo, err: err}:
		case <-done:
			return
		}
	}

}

func loopThroughRepo(done <-chan struct{}, repoList []RepositoryGithubDto) <-chan RepositoryGithubDto {
	out := make(chan RepositoryGithubDto)

	go func() {
		defer close(out)
		for _, repoDto := range repoList {
			out <- repoDto

			// Abort listing repositories if done is closed.
			select {
			case <-done:
				fmt.Println("listing repositories canceled")
				return
			default:
				break
			}

		}
	}()
	return out
}

func merge(done <-chan struct{}, cs []<-chan Repository) <-chan Repository {
	var wg sync.WaitGroup
	out := make(chan Repository, 1) // enough space for the unread inputs

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c or done is closed, then calls wg.Done.
	output := func(c <-chan Repository) {
		defer wg.Done()
		for n := range c {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func listRepoToTheChan(done <-chan struct{}, repoList []RepositoryGithubDto) <-chan RepositoryGithubDto {
	out := make(chan RepositoryGithubDto, len(repoList))
	go func() {
		defer close(out)
		for _, repo := range repoList {
			select {
			case out <- repo:
			case <-done:
				return
			}
		}
		close(out)
	}()
	return out
}

func fetchAndTransferRepoToTheChan(done <-chan struct{}, in <-chan RepositoryGithubDto) <-chan Repository {
	out := make(chan Repository)

	go func() {
		for repoDto := range in {
			languages, err := FetchLanguagesList(repoDto.Owner.OwnerName, repoDto.Name)
			if err != nil {
				return
			}
			repo := Repository{
				Name:        repoDto.Name,
				FullName:    repoDto.FullName,
				Description: repoDto.Description,
				OwnerName:   repoDto.Owner.OwnerName,
				AvatarURL:   repoDto.Owner.AvatarURL,
				Type:        repoDto.Owner.Type,
				URL:         repoDto.URL,
				Languages:   languages,
			}
			select {
			case out <- repo:
			case <-done:
			}

		}
		close(out)
	}()
	return out
}
