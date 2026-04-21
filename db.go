package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./site.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	// Create tables if they don't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS articles (
			slug TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS views (
			slug TEXT PRIMARY KEY,
			count INTEGER DEFAULT 0
		);
	`)
	if err != nil {
		log.Fatal("Failed to create tables:", err)
	}

	// Seed a sample article if none exist
	var count int
	db.QueryRow("SELECT COUNT(*) FROM articles").Scan(&count)
	if count == 0 {
		db.Exec(`
			INSERT INTO articles (slug, title, content) VALUES 
			('first-post', 'My First Article', 'This is the content of my first article. Welcome to my site!')
		`)
	}

	log.Println("Database ready")
}

func incrementView(slug string) {
	db.Exec(`
		INSERT INTO views (slug, count) VALUES (?, 1)
		ON CONFLICT(slug) DO UPDATE SET count = count + 1
	`, slug)
}

func getViewCount(slug string) int {
	var count int
	db.QueryRow("SELECT count FROM views WHERE slug = ?", slug).Scan(&count)
	return count
}

func getArticle(slug string) (Article, error) {
	var a Article
	err := db.QueryRow(
		"SELECT slug, title, content, created_at FROM articles WHERE slug = ?", slug,
	).Scan(&a.Slug, &a.Title, &a.Content, &a.CreatedAt)
	return a, err
}

func getAllArticles() ([]Article, error) {
	rows, err := db.Query("SELECT slug, title, created_at FROM articles ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var a Article
		rows.Scan(&a.Slug, &a.Title, &a.CreatedAt)
		articles = append(articles, a)
	}
	return articles, nil
}
