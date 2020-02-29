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

type languageStats struct {
	count    uint16
	size     uint64
	repoList map[string]Repository // map[fullname: string]URL: string
}

// Pipeline:
// 2 stages
// 1st: list repo to chan
// 2nd:
//		- fetch languages
// 		- aggregate repo: map[language: string]{count: number, size: number}

// getAggregatedRepo will fetch public repositories, their languages then call worker to aggregate those data
func getAggregatedRepo() (map[string]languageStats, error) {
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
	numWorker := serverConfig.WorkerNumber
	wg.Add(numWorker)
	for i := 0; i < numWorker; i++ {
		go func() {
			fetchLanguagesWorker(done, repoCh, c)
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(c)
	}()
	// End of pipeline.

	languagesStats := aggregateData(c)

	elapsed := time.Since(start)

	fmt.Println("\nfinal stats:")

	for lang, stats := range languagesStats {
		fmt.Println("\n	->", lang, "  count:", stats.count, "  size:", stats.size)
		for name := range stats.repoList {
			fmt.Println("		-", name)
		}
	}

	fmt.Printf("Fetch repositories with their languages list took %s", elapsed)
	return languagesStats, nil
}

func aggregateData(c chan result) map[string]languageStats {
	// RepoList use a map where key is full_name field.
	// full_name field is unique for github
	// repoList := make(map[string]Repository)
	aggregatedData := make(map[string]languageStats)

	for r := range c {
		if r.err != nil {
			fmt.Println("Error detect during fetching languages from repository", r.err)
			continue
		}
		// repoList[r.repository.FullName] = r.repository
		for lang, size := range r.repository.Languages {
			var tmpRepoList map[string]Repository

			if aggregatedData[lang].repoList == nil {
				tmpRepoList = make(map[string]Repository)
			} else {
				tmpRepoList = aggregatedData[lang].repoList
			}
			tmpStats := &languageStats{
				count:    aggregatedData[lang].count + 1,
				size:     aggregatedData[lang].size + uint64(size),
				repoList: tmpRepoList,
			}

			aggregatedData[lang] = *tmpStats
			aggregatedData[lang].repoList[r.repository.FullName] = r.repository
		}
		// fmt.Println(r.repository)
	}

	return aggregatedData
}

func fetchLanguagesWorker(done <-chan struct{}, repositoriesDto <-chan RepositoryGithubDto, c chan<- result) {
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
