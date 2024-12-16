package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
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
