package common

import (
	"context"
	"encoding/json"
	"log/slog"
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
	ctx := context.TODO()
	tmp, err := json.Marshal(t)
	if err != nil {
		slog.ErrorContext(ctx,
			"Message Json error",
			"Error", err,
		)
		return nil
	}
	return tmp
}
