package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

func main() {
	url := "postgres://postgres:mysecretpassword@localhost:5432/wiki-db"
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), "select title from article")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
		os.Exit(1)
	}

	titles := make([]string, 0)
	for rows.Next() {
		title := ""
		rows.Scan(&title)
		titles = append(titles, title)
	}

	speedTest(titles)
}

func speedTest(titles []string) {
	const num_queries = 1_000_000
	const num_routines = 95

	channels := make([]chan string, 0, num_routines)
	returnChannels := make([]chan bool, 0, num_routines)
	for i := 0; i < num_routines; i++ {
		channels = append(channels, make(chan string, num_queries/num_routines))
		returnChannels = append(returnChannels, make(chan bool, 0))
		go query(channels[i], returnChannels[i])
	}

	start_time := time.Now()
	for i := 0; i < num_queries; i++ {
		index := rand.Int() % len(titles)
		channels[i%num_routines] <- titles[index]
	}
	for i, _ := range channels {
		channels[i] <- ""
		close(channels[i])
	}

	for i, _ := range returnChannels {
		hello := <-returnChannels[i]
		fmt.Println(hello)
	}
	end_time := time.Now()

	fmt.Printf("Time Elapsed: %d\n", end_time.Unix()-start_time.Unix())
	fmt.Printf("Queries Per Second: %f\n", num_queries/float64(end_time.Unix()-start_time.Unix()))
}

func query(ch chan string, r chan bool) {
	url := "postgres://postgres:mysecretpassword@localhost:5432/wiki-db"
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	for title := "a"; title != ""; title = <-ch {
		title := <-ch
		conn.QueryRow(context.Background(), "select * from article where title = '"+title+"'").Scan(nil)
	}
	r <- true
	close(r)
}
