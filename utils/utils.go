package utils

import (
	"encoding/json"
	"io/ioutil"
)

// FileData ..
type FileData struct {
	Data interface{}
}

// ReadFile ..
func ReadFile(filename string, fileContent FileData) error {
	byt, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(byt, fileContent.Data)
}
