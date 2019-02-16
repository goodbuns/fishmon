// Package adafruitio provides an API client for Adafruit.IO.
package adafruitio

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

// A Client provides authentication for Adafruit.IO API requests.
type Client struct {
	username string
	aioKey   string
}

// New constructs an Adafruit.IO API client and tests its validity.
func New(username, AIOKey string) (*Client, error) {
	// Construct client.
	client := &Client{
		username: username,
		aioKey:   AIOKey,
	}

	// Check client credentials.
	req, err := http.NewRequest(http.MethodGet, "https://io.adafruit.com/api/v2/user", nil)
	if err != nil {
		return nil, err
	}
	res, err := client.send(req)
	if err != nil {
		return nil, err
	}
	err = check(res)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Record uploads a value to an Adafruit.IO feed.
func (c *Client) Record(feed, value string, timestamp time.Time) error {
	// Marshal request body.
	payload, err := json.Marshal(request{
		Value:     value,
		CreatedAt: timestamp,
	})
	if err != nil {
		return err
	}

	// Set request headers.
	req, err := http.NewRequest(
		http.MethodPost,
		"https://io.adafruit.com/api/v2/"+c.username+"/feeds/"+feed+"/data",
		bytes.NewReader(payload),
	)
	if err != nil {
		return err
	}

	// Send request.
	res, err := c.send(req)
	if err != nil {
		return err
	}

	// Check response for errors.
	defer res.Body.Close()
	return check(res)
}

func (c *Client) send(req *http.Request) (*http.Response, error) {
	req.Header.Set("X-AIO-Key", c.aioKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return http.DefaultClient.Do(req)
}

func check(res *http.Response) error {
	// Parse response body.
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var r response
	err = json.Unmarshal(body, &r)
	if err != nil {
		return err
	}

	// Check for application-level errors.
	if r.Error != "" {
		return errors.New(r.Error)
	}

	return nil
}

type request struct {
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
}

type response struct {
	Error string
}
