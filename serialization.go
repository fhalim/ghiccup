package ghiccup

import (
	"encoding/json"
	"log"
)

func Marshall(value interface{}) string {
	bytes, err := json.Marshal(value)
	if err != nil {
		log.Panicf("Error %v", err)
	}
	return string(bytes)
}
