package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/coderjz/gophercises/exercise5/sitemap"
)

func main() {
	var startURL string
	var maxDepth int
	flag.StringVar(&startURL, "url", "https://gophercises.com/", "URL to begin our sitemap")
	flag.IntVar(&maxDepth, "depth", 3, "Maximum depth for our search for URLs in the sitemap")
	flag.Parse()

	allLinks := sitemap.FindAllLinksOnDomain(startURL, maxDepth)

	sitemapStr, err := sitemap.GenerateSitemap(allLinks)
	if err != nil {
		log.Fatalf("Error generating sitemap - %v\n", err)
	}
	fmt.Println(sitemapStr)
}
