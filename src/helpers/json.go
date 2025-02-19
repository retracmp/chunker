package helpers

import (
	"os"

	"github.com/goccy/go-json"
)

type JSON = map[string]interface{}

func ToJSONString(v interface{}) string {
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(out)
}

func ToJSON[T any](v T) JSON {
	var j JSON
	data, err := json.Marshal(v)
	if err != nil {
		return j
	}

	if err := json.Unmarshal(data, &j); err != nil {
		return j
	}

	return j
}

func JSONFromString[T any](input string) T {
	var output T
	json.Unmarshal([]byte(input), &output)
	return output
}

func JSONFromBytes[T any](input []byte) T {
	var output T
	if err := json.Unmarshal(input, &output); err != nil {
		panic(err)
	}
	return output
}

func JSONToBytes[T any](v T) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}

func JSONFromFile[T any](path string) T {
	var output T
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(data, &output); err != nil {
		panic(err)
	}

	return output
}