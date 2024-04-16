package common

import (
	"encoding/json"
	"log"
)

type message struct {
	Results interface{} `json:"Result"`
}

func Message(str interface{}) message {
	return message{
		Results: str,
	}
}

func (t *message) Json() []byte {
	tmp, err := json.Marshal(t)
	if err != nil {
		log.Println("error:", err)
		return nil
	}
	return tmp
}
