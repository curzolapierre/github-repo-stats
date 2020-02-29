package main

import (
	"log"
	"net/http"
	"regexp"
	"text/template"
)

var validPath = regexp.MustCompile("^/|(index|search)|/([a-zA-Z0-9]+)$")
var templates = template.Must(template.ParseFiles(
	"template/index.html",
	"template/search.html"))

type inputPage struct {
	RepoName     []byte
	languageName []byte
}

// Page structure handle variables sent to client
type Page struct {
	Body []byte
}

func executeSearch(repoName, languageName string) *map[string]languageStats {
	var params string

	if repoName != "" {
		params = "q=in:name+" + repoName + "+"
	}
	if languageName != "" {
		if params == "" {
			params = "q=language:" + languageName
		} else {
			params += "language:" + languageName
		}
	}
	repoStats, err := getAggregatedRepo(params)
	if err != nil {
		log.Print(err)
	}
	return &repoStats
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

func renderTemplate(w http.ResponseWriter, template string, p interface{}) {
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
	languageName := r.FormValue("languageName")
	repo := executeSearch(repoName, languageName)

	renderTemplate(w, "search", repo)
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{}
	renderTemplate(w, "error", p)
}
