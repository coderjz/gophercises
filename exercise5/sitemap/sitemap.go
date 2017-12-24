package sitemap

import (
	"encoding/xml"
)

//GenerateSitemap will create the actual sitemap string
func GenerateSitemap(urls []string) (string, error) {
	urlset := urlset{
		Namespace: "http://www.sitemaps.org/schemas/sitemap/0.9",
		Url:       []urlEntry{},
	}

	for _, url := range urls {
		urlset.Url = append(urlset.Url, urlEntry{
			Url: url,
		})
	}

	result, err := xml.MarshalIndent(urlset, "", "  ")
	if err != nil {
		return "", err
	}
	return xml.Header + string(result), nil
}

type urlset struct {
	Namespace string     `xml:"xmlns,attr"`
	Url       []urlEntry `xml:"url"`
}

type urlEntry struct {
	Url string `xml:"loc"`
}
