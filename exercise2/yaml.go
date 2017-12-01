package exercise2

import (
	yaml "gopkg.in/yaml.v2"
)

func parseYAML(yml []byte) ([]map[string]string, error) {
	urlList := make([]map[string]string, 0)

	err := yaml.Unmarshal(yml, &urlList)
	if err != nil {
		return nil, err
	}
	return urlList, nil
}

func buildMap(yamlMapList []map[string]string) map[string]string {
	mergedMap := make(map[string]string)
	for _, yamlMap := range yamlMapList {
		mergedMap[yamlMap["path"]] = yamlMap["url"]
	}
	return mergedMap
}
