package utils

import "encoding/json"

func ToJSON(object interface{}) string {
	b, err := json.Marshal(&object)
	if err != nil {
		panic(err)
	}
	return string(b)
}
