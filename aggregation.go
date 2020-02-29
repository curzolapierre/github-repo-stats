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
	Count    uint16
	Size     uint64
	RepoList map[string]Repository
}

func aggregateData(c chan result) map[string]languageStats {
	aggregatedData := make(map[string]languageStats)

	for r := range c {
		if r.err != nil {
			fmt.Println("Error detect during fetching languages from repository", r.err)
			continue
		}
		for lang, size := range r.repository.Languages {
			// tmpRepoList use a map where key is unique 'full_name' field.
			var tmpRepoList map[string]Repository

			if aggregatedData[lang].RepoList == nil {
				tmpRepoList = make(map[string]Repository)
			} else {
				tmpRepoList = aggregatedData[lang].RepoList
			}
			tmpStats := &languageStats{
				Count:    aggregatedData[lang].Count + 1,
				Size:     aggregatedData[lang].Size + uint64(size),
				RepoList: tmpRepoList,
			}

			aggregatedData[lang] = *tmpStats
			aggregatedData[lang].RepoList[r.repository.FullName] = r.repository
		}
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

// Pipeline:
// 2 stages
// 1st: list repo to chan
// 2nd:
//		- fetch languages
// 		- aggregate repo: map[language: string]{count: number, size: number, map: Repository}

// getAggregatedRepo will fetch public repositories, their languages then call worker to aggregate those data
func getAggregatedRepo(querySearch ...string) (map[string]languageStats, error) {
	start := time.Now()

	repoDtoList, err := FetchRepositoriesList(querySearch...)
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

	fmt.Printf("Fetch repositories with their languages list took %s", elapsed)
	return languagesStats, nil
}
