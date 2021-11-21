package helper

import (
	"bytes"
	"encoding/json"
)

func StructToJSONReader(data interface{}) *bytes.Reader {
	userJSON, _ := json.Marshal(data)
	return bytes.NewReader(userJSON)
}
