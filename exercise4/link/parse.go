package link

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

//Link represents a single link in an HTML document (<a href="...">...</a>)
type Link struct {
	Href string
	Text string
}

//Parse will take in an HTML document and will return a slice of links found
func Parse(reader io.Reader) ([]Link, error) {
	doc, err := html.Parse(reader)
	if err != nil {
		return nil, fmt.Errorf("Error parsing html: %v", err)
	}
	links := make([]Link, 0)
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			links = append(links, buildLink(n))
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return links, nil
}

func buildLink(n *html.Node) Link {
	href := ""
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			href = attr.Val
		}
		break
	}

	var text func(*html.Node) string
	text = func(n *html.Node) string {
		if n.Type == html.TextNode {
			return n.Data
		}
		if n.Type != html.ElementNode {
			return ""
		}
		var ret string
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			ret += text(c)
		}
		return strings.Join(strings.Fields(ret), " ")
	}
	linkText := strings.TrimSpace(text(n))

	return Link{
		Href: href,
		Text: linkText,
	}
}
