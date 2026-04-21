package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type PageData struct {
	Title            string
	Articles         []Article
	Article          Article
	SubscribeSuccess bool
	SubscribeError   string
}

var templates = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	initDB()

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/article/", articleHandler)
	http.HandleFunc("/subscribe", subscribeHandler)

	port := "8080"
	fmt.Println("Server running on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	articles, err := getAllArticles()
	if err != nil {
		http.Error(w, "Could not load articles", http.StatusInternalServerError)
		return
	}

	data := PageData{
		Title:    "Home | My Site",
		Articles: articles,
	}
	err = templates.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func articleHandler(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/article/")
	if slug == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	article, err := getArticle(slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Increment and fetch view count
	incrementView(slug)
	article.Views = getViewCount(slug)

	data := PageData{
		Title:   article.Title + " | My Site",
		Article: article,
	}
	err = templates.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}


func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	email := r.FormValue("email")
	if email == "" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	articles, _ := getAllArticles()
	data := PageData{
		Title:    "Home | My Site",
		Articles: articles,
	}

	err := addContactToBrevo(email)
	if err != nil {
		data.SubscribeError = "Something went wrong. Please try again."
	} else {
		data.SubscribeSuccess = true
	}

	templates.ExecuteTemplate(w, "base.html", data)
}



