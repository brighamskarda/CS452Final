package wikicalls

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"

	"golang.org/x/net/html"
)

type Article struct {
	Title string
	Links []string
}

func GetArticleWiki(title string) (Article, error) {
	request, err := http.NewRequest("GET", "https://en.wikipedia.org/api/rest_v1/page/html/"+title, nil)
	if err != nil {
		fmt.Fprint(os.Stderr, "Error forming request in GetArticleWiki\n")
		return Article{}, err
	}
	request.Header.Add("User-Agent", "brighamskarda@gmail.com")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return Article{}, err
	}

	article_title := getArticleTitle(response)
	links := getArticleLinks(response.Body)

	return Article{article_title, links}, err
}

func getArticleTitle(response *http.Response) string {
	content_location := response.Header.Get("content-location")
	if content_location == "" {
		return ""
	}
	split_content_location := strings.Split(content_location, "/")
	return split_content_location[len(split_content_location)-1]
}

func getArticleLinks(body io.ReadCloser) []string {
	node, err := html.Parse(body)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to parse HTML in getArticleLinks")
		return nil
	}

	var links []string
	for descendant := range node.Descendants() {
		link := getLink(descendant)
		if link != "" {
			links = append(links, link)
		}
	}

	return links
}

func getLink(node *html.Node) string {
	if slices.ContainsFunc(node.Attr, func(a html.Attribute) bool {
		return a.Key == "rel" && a.Val == "mw:WikiLink"
	}) {
		for _, a := range node.Attr {
			if a.Key == "href" && !strings.ContainsRune(a.Val, '#') && strings.Count(a.Val, "/") == 1 {
				return a.Val[2:]
			}
		}
	}
	return ""
}
