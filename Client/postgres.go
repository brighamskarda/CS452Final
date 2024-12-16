package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func GetArticlesPostgres(titles []string) []Article {
	url := "postgres://postgres:mysecretpassword@localhost:5432/wiki-db"
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	const query = "SELECT * FROM article WHERE title = $1"

	articles := make([]Article, len(titles))
	for i, title := range titles {
		row := conn.QueryRow(context.Background(), query, title)
		articles[i].Title = title
		if row.Scan(articles[i].Links) != nil {
			articles[i].Title = ""
			articles[i].Links = nil
		}
	}
	return articles
}

func InsertArticlesPostgres(articles []Article) {
	url := "postgres://postgres:mysecretpassword@localhost:5432/wiki-db"
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	const query = "INSERT INTO article VALUES($1, $2)"

	for _, article := range articles {
		commandTag, err := conn.Exec(context.Background(), query, article.Title, article.Links)
		if err != nil || commandTag.RowsAffected() != 1 {
			fmt.Fprintln(os.Stderr, "Error inserting row in PostArticlesPostgres.")
		}
	}
}