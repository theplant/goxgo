package goxgo

import (
	"encoding/json"
	"fmt"
)

// Serialization function for the request structures into JSON
func Serialize(i interface{}) (payload []byte, err error) {
	payload, err = json.Marshal(i)
	if err != nil {
		fmt.Println("error:", err)
	}
	return
}

// Unserialization function for the structures from JSON to a
// response structure
func Unserialize(payload []byte, i interface{}) (err error) {
	err = json.Unmarshal(payload, &i)
	if err != nil {
		fmt.Printf("error: %s\n%s\n", err, payload)
	}
	return
}
