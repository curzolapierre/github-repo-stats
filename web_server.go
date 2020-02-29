package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"text/template"
)

var validPath = regexp.MustCompile("^/|(index|error|search)|/([a-zA-Z0-9]+)$")
var templates = template.Must(template.ParseFiles(
	"./template/index.html",
	"./template/search.html",
	"./template/error.html"))

// Page structure handle variables sent to client
type Page struct {
	Body []byte
}

func (p *Page) executeSearch() {
	if p.Body != nil && string(p.Body) != "" {
		repoName := string(p.Body)
		fmt.Println("query search:", repoName)

		params := "q=" + repoName
		repoStats, err := getAggregatedRepo(params)
		if err != nil {
			log.Fatalln(err)
		}

		var jsonData []byte
		jsonData, err = json.Marshal(repoStats)
		if err != nil {
			log.Println(err)
		}
		p.Body = jsonData
	}
}

func makeHandler(httpFunction func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		httpFunction(w, r)
	}
}

func renderTemplate(w http.ResponseWriter, template string, p *Page) {
	err := templates.ExecuteTemplate(w, template+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{}
	renderTemplate(w, "index", p)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	repoName := r.FormValue("repoName")
	p := &Page{Body: []byte(repoName)}
	p.executeSearch()
	// err := p.executeSearch()
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// http.Redirect(w, r, "/search/", http.StatusFound)
	w.Header().Set("Content-Type", "application/json")
	// jsonObj, err := json.Marshal(p)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// }

	w.Write(p.Body)
	// renderTemplate(w, "search", p)
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{}
	renderTemplate(w, "error", p)
}
