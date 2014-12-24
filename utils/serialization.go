package utils

import (
	"encoding/json"
	"fmt"
)

func Marshall(value interface{}) string {
	bytes, err := json.Marshal(value)
	if err != nil {
		panic(fmt.Sprintf("Error %v", err))
	}
	return string(bytes)
}
