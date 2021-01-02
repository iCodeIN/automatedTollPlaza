package utils

import (
	"encoding/json"
	"fmt"
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

// PrintJSON ..
func PrintJSON(v interface{}) {
	byt, _ := json.Marshal(v)
	fmt.Println(string(byt))
}
