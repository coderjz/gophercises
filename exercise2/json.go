package exercise2

import (
	yaml "gopkg.in/yaml.v2"
)

type jsonURLList struct {
	Urls []jsonURL `json:"urls"`
}

type jsonURL struct {
	URL  string `json:"url"`
	Path string `json:"path"`
}

func parseJSON(json []byte) (jsonURLList, error) {
	var urlList jsonURLList

	err := yaml.Unmarshal(json, &urlList)
	return urlList, err
}

func buildJSONMap(urllist jsonURLList) map[string]string {
	mergedMap := make(map[string]string)
	for _, url := range urllist.Urls {
		mergedMap[url.Path] = url.URL
	}
	return mergedMap
}
