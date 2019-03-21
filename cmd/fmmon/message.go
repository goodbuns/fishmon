package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

type Message struct {
	Text string `json:"text"`
}

func Send(url, msg string) error {
	m := Message{
		Text: msg,
	}

	body, err := json.Marshal(m)
	if err != nil {
		return errors.Wrap(err, "could not marshal alert message")
	}

	_, err = http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return errors.Wrap(err, "could not send alert message")
	}

	return nil
}
