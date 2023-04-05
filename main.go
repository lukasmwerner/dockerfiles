package main

import (
	"embed"
	"errors"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"

	_ "embed"

	"github.com/lukasmwerner/pine"
	"github.com/russross/blackfriday/v2"
)

//go:embed templates/lang.tmpl.html
var postPage string

//go:embed templates/home.tmpl.html
var homePage string

//go:embed languages
var postsFS embed.FS

type Post struct {
	Title   string
	Content string
}

func main() {
	postTemplate := template.Must(template.New("post").Parse(postPage))
	homeTemplate := template.Must(template.New("home").Parse(homePage))

	posts := make(map[string]Post)
	results, _ := fs.Glob(postsFS, "**/*.md")
	for _, result := range results {

		b, err := fs.ReadFile(postsFS, result)
		if err != nil {
			panic(err)
		}
		parsed := blackfriday.Run(b, blackfriday.WithExtensions(blackfriday.CommonExtensions))

		url := "/" + strings.Replace(strings.ReplaceAll(result, ".md", ""), "languages", "lang", 1)
		posts[url] = Post{
			Title:   url,
			Content: string(parsed),
		}
	}

	p := pine.New()
	fs := http.FileServer(http.Dir("./static"))
	p.Handle("/static", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/static", fs).ServeHTTP(w, r)
	})
	p.Handle("/lang/{key}", func(w http.ResponseWriter, r *http.Request) {

		p, ok := posts[r.URL.Path]
		if !ok {
			http.Error(w, errors.New("post not found").Error(), 404)
			return
		}

		postTemplate.Execute(w, p)
	})
	p.Handle("/", func(w http.ResponseWriter, r *http.Request) {
		links := make([]string, len(posts))
		i := 0
		for k := range posts {
			links[i] = k
			i++
		}
		homeTemplate.Execute(w, struct {
			Posts []string
		}{
			Posts: links,
		})
	})

	srv := &http.Server{
		Handler: p,
		Addr:    "0.0.0.0:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Listening on 0.0.0.0:8000")
	log.Fatal(srv.ListenAndServe())
}
