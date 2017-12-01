package exercise2

import (
	"fmt"
	"net/http"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		redirectURL, ok := pathsToUrls[request.URL.Path]
		if ok {
			//Do redirect here
			http.Redirect(writer, request, redirectURL, http.StatusPermanentRedirect)
			return
		}
		fallback.ServeHTTP(writer, request)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yaml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yaml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

//JSONHandler is the bonus exercise for parsing JSON with path value
func JSONHandler(json []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJSON, err := parseJSON(json)
	if err != nil {
		return nil, err
	}
	pathMap := buildJSONMap(parsedJSON)
	return MapHandler(pathMap, fallback), nil
}

//BoltDBHandler is the bonus exercise for parsing from boltDB
func BoltDBHandler(fallback http.Handler) (http.HandlerFunc, error) {
	pathMap, err := getBoltRecords()
	if err != nil {
		return nil, err
	}
	fmt.Println(pathMap)
	return MapHandler(pathMap, fallback), nil
}
