package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"math"
	"os"
	"strings"
)

func main() {
	setLogLevel()
	fmt.Print("Welcome to wiki path finder! Enter the names of any two " +
		"wikipedia articles and this program will attempt to find a series of links that connect them.\n\n")

	for {
		fromArticleTitle, toArticleTitle := promptUserForArticles()

		fromArticle, err := GetArticleWiki(fromArticleTitle)
		if err != nil {
			fmt.Println("Couldn't find article titled ", fromArticleTitle)
			continue
		}
		toArticle, err := GetArticleWiki(toArticleTitle)
		if err != nil {
			fmt.Println("Couldn't find article titled ", fromArticleTitle)
			continue
		}

		if cachedValue := CheckRedisCache(fromArticle.Title, toArticle.Title); cachedValue != "" {
			fmt.Print("Found cached path: \033[1m" + cachedValue + "\033[0m\n\n")
			continue
		}

		path := BreadthFirstSearch(fromArticle, toArticle)
		if path == "" {
			fmt.Print("Could not find path from ", fromArticle.Title, " to ", toArticle.Title, "\n\n")
			continue
		}

		fmt.Print("Found path: \033[1m" + path + "\033[0m\n\n")
		InsertRedisCache(fromArticle.Title+"->"+toArticle.Title, path)
	}
}

func setLogLevel() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "DEBUG":
			slog.SetLogLoggerLevel(slog.LevelDebug)
		case "INFO":
			slog.SetLogLoggerLevel(slog.LevelInfo)
		case "WARN":
			slog.SetLogLoggerLevel(slog.LevelWarn)
		case "ERROR":
			slog.SetLogLoggerLevel(slog.LevelError)
		default:
			slog.SetLogLoggerLevel(math.MinInt)
		}
	}
}

func promptUserForArticles() (string, string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter starting article title (spaces allowed):")
	fromArticle, _ := reader.ReadString('\n')
	fromArticle = strings.TrimSpace(fromArticle)
	fmt.Println("Enter destination article title (spaces allowed):")
	toArticle, _ := reader.ReadString('\n')
	toArticle = strings.TrimSpace(toArticle)
	return fromArticle, toArticle
}

type BfsNode struct {
	title    string
	children []*BfsNode
}

const bfsMaxDepth = 6

func BreadthFirstSearch(fromArticle Article, toArticle Article) string {
	requestChannel := make(chan apiRequest, apiRequestSec)
	go ApiRequestProcessor(requestChannel)

	parentNode := &BfsNode{
		title:    fromArticle.Title,
		children: make([]*BfsNode, 0, len(fromArticle.Links)),
	}
	// fill in parent node children
	for _, link := range fromArticle.Links {
		parentNode.children = append(parentNode.children, &BfsNode{
			title:    link,
			children: nil})
		if link == toArticle.Title {
			return fromArticle.Title + "->" + toArticle.Title
		}
	}

	for depthToSearch := 0; depthToSearch <= bfsMaxDepth; depthToSearch++ {
		if result := search(parentNode, requestChannel, depthToSearch, toArticle.Title); result != "" {
			if depthToSearch == 0 {
				return parentNode.title + "->" + result
			}
			return result
		}
	}

	return ""
}

func search(node *BfsNode, requestChan chan apiRequest, depth int, toArticle string) string {
	// recurse
	if depth != 0 {
		for _, child := range node.children {
			if result := search(child, requestChan, depth-1, toArticle); result != "" {
				return node.title + "->" + child.title + "->" + result
			}
		}
		return ""
	}

	// base case
	child_titles := make([]string, 0, len(node.children))
	for _, child := range node.children {
		child_titles = append(child_titles, child.title)
	}
	articles := GetArticles(child_titles, requestChan)

	for i, child := range node.children {
		grandChildren := make([]*BfsNode, 0, len(articles[i].Links))
		for _, link := range articles[i].Links {
			grandChildren = append(grandChildren, &BfsNode{title: link, children: nil})
			if link == toArticle {
				return child.title + "->" + link
			}
		}
		child.children = grandChildren
	}

	return ""
}

func GetArticles(titles []string, requestChan chan apiRequest) []Article {
	articles := GetArticlesPostgres(titles)

	returnChannels := make([]chan Article, 0)
	for i, article := range articles {
		if article.Title == "" {
			newChan := make(chan Article)
			request := apiRequest{
				title:         titles[i],
				returnChannel: newChan,
			}
			returnChannels = append(returnChannels, newChan)
			requestChan <- request
		}
	}

	articlesToInsert := []Article{}
	for i, article := range articles {
		if article.Title == "" {
			articles[i] = <-returnChannels[0]
			returnChannels = returnChannels[1:]
			articlesToInsert = append(articlesToInsert, articles[i])
		}
	}

	InsertArticlesPostgres(articlesToInsert)

	return articles
}
