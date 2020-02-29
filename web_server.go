package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"text/template"
)

var validPath = regexp.MustCompile("^/|(index|error|search)|/([a-zA-Z0-9]+)$")
var templates = template.Must(template.ParseFiles(
	"template/index.html",
	"template/search.html",
	"template/error.html"))

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
	query := r.FormValue("repoName")
	p := &Page{Body: []byte(query)}
	// err := p.executeSearch()
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// http.Redirect(w, r, "/search/", http.StatusFound)
	// w.Header().Set("Content-Type", "application/json")

	// w.Write(p.Body)
	renderTemplate(w, "search", p)
}

func queryHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.Write([]byte("nothing to do"))
		return
	}
	type bodyReader struct {
		QuerySearch string `json:"querySearch"`
	}
	body, err := ioutil.ReadAll(r.Body)
	fmt.Println("data receive from request", string(body))
	data := bodyReader{}

	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = json.Unmarshal([]byte(body), &data)

	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Println("data receive from request", data)
	p := &Page{Body: []byte(data.QuerySearch)}
	p.executeSearch()

	w.Header().Set("Content-Type", "application/json")

	w.Write(p.Body)
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{}
	renderTemplate(w, "error", p)
}
