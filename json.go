package json2wasi2json

import "encoding/json"

// Bytes2JsonMap converts a byte slice into a JsonMap.
func Bytes2JsonMap(data []byte) (JsonMap, error) {
	var m JsonMap
	err := json.Unmarshal(data, &m)
	return m, err
}

// JsonMap2Bytes converts a JsonMap into a byte slice.
func JsonMap2Bytes(m JsonMap) ([]byte, error) {
	return json.Marshal(m)
}
