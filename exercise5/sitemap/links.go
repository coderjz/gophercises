package sitemap

import (
	"crypto/tls"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/coderjz/gophercises/exercise4/link"
)

//FindAllLinksOnDomain returns all found URLs that can be found that are pointing to the same domain
func FindAllLinksOnDomain(startURL string, maxDepth int) []string {
	currentDepth := 1
	currentDepthEnd := 0

	domain := getDomain(startURL)
	if domain == "" {
		log.Fatalf("Invalid starting URL %v - could not find domain.\n", startURL)
	}

	startURL = canonicalizeURL(startURL, domain, "")
	allURLs := []string{startURL}
	for i := 0; i < len(allURLs); i++ {
		if i >= currentDepthEnd {
			currentDepth++
			currentDepthEnd = len(allURLs)
			if currentDepth > maxDepth {
				break
			}
		}

		//Warning: Gophercises video mentions that a better solution is to use the final (redirected) URL
		//obtained from resp.Request.URL from our GET request.
		//One approach could be to return two values from this function and then update allURLs[i], but I really don't like the double return value.
		//So leave-as is for now.
		nextLinks := getLinksFromURL(allURLs[i])

		for _, link := range nextLinks {
			if !isValidLink(link, domain) {
				continue
			}

			canonicalLinkURL := canonicalizeURL(link.Href, domain, allURLs[i])
			linkExists := false
			for _, existingLink := range allURLs {
				if canonicalLinkURL == existingLink {
					linkExists = true
					break
				}
			}
			if !linkExists {
				allURLs = append(allURLs, canonicalLinkURL)
			}
		}
	}
	return allURLs
}

func getLinksFromURL(url string) []link.Link {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		return []link.Link{}
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return []link.Link{}
	}
	nextLinks, err := link.Parse(resp.Body)
	if err != nil {
		return []link.Link{}
	}
	return nextLinks
}

func getDomain(startURL string) string {
	re := regexp.MustCompile(`^\s*(?:https?://)?([^/]+)(?:/|$)`)
	matches := re.FindAllStringSubmatch(startURL, 1)
	if len(matches) == 0 {
		return ""
	}
	return matches[0][1]
}

func isValidLink(link link.Link, domain string) bool {
	if strings.HasPrefix(strings.TrimSpace(link.Href), "/") {
		return true
	}

	if getDomain(strings.TrimSpace(link.Href)) == domain {
		return true
	}
	return false
}

func canonicalizeURL(url string, domain string, parentURL string) string {
	url = strings.TrimSpace(url)
	parentPageProtocol := getProtocol(parentURL)
	if strings.HasPrefix(url, "//") {
		return parentPageProtocol + ":" + url
	}
	if strings.HasPrefix(url, "/") {
		return parentPageProtocol + "//" + strings.TrimSuffix(domain, "/") + url
	}

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	if url == "http://"+domain || url == "https://"+domain {
		url = url + "/"
	}
	return url
}

func getProtocol(url string) string {
	if strings.HasPrefix(url, "https://") {
		return "https"
	}
	return "http"
}
