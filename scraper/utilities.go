package scraper

import (
	"bytes"
	"encoding/json"
)

func mergeResults(dest results, src results) {
	for k, v := range src {
		for _, text := range v {
			dest[k] = append(dest[k], text)
		}
	}
}

// JSONMarshalUnescaped does what json.Marshal does without escaping &, <, >.
func JSONMarshalUnescaped(v interface{}) ([]byte, error) {
	data, err := json.Marshal(v)
	data = bytes.Replace(data, []byte("\\u0026"), []byte("&"), -1)
	data = bytes.Replace(data, []byte("\\u003c"), []byte("<"), -1)
	data = bytes.Replace(data, []byte("\\u003e"), []byte(">"), -1)
	return data, err
}
