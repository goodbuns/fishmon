// Package adafruitio provides an API client for Adafruit.IO.
package adafruitio

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
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
		return nil, errors.Wrap(err, "could not construct API request")
	}
	res, err := client.send(req)
	if err != nil {
		return nil, errors.Wrap(err, "could not send API request")
	}
	err = check(res)
	if err != nil {
		return nil, errors.Wrap(err, "API response has error")
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
		return errors.Wrap(err, "could not marshal API request body")
	}

	// Set request headers.
	req, err := http.NewRequest(
		http.MethodPost,
		"https://io.adafruit.com/api/v2/"+c.username+"/feeds/"+feed+"/data",
		bytes.NewReader(payload),
	)
	if err != nil {
		return errors.Wrap(err, "could not set API request headers")
	}

	// Send request.
	res, err := c.send(req)
	if err != nil {
		return errors.Wrap(err, "could not send API request")
	}

	// Check response for errors.
	defer res.Body.Close()
	err = check(res)
	if err != nil {
		return errors.Wrap(err, "API response has errors")
	}
	return nil
}

// GetFeed returns all feeds of the Adafruit user.
func (c *Client) GetFeed() ([]Feed, error) {
	// Create request.
	req, err := http.NewRequest(
		http.MethodGet,
		"https://io.adafruit.com/api/v2/"+c.username+"/feeds",
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "could not set API request headers")
	}

	// Send request.
	res, err := c.send(req)
	if err != nil {
		return nil, errors.Wrap(err, "could not send API request")
	}
	defer res.Body.Close()

	// Marshal response into Feed struct.
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "could not read response body")
	}

	var feeds []Feed
	err = json.Unmarshal(body, &feeds)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshal response body for feed data")
	}

	return feeds, nil
}

func (c *Client) send(req *http.Request) (*http.Response, error) {
	req.Header.Set("X-AIO-Key", c.aioKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "could not send HTTP request")
	}
	return res, nil
}

func check(res *http.Response) error {
	// Parse response body.
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "could not read response body")
	}
	var r response
	err = json.Unmarshal(body, &r)
	if err != nil {
		return errors.Wrap(err, "could not unmarshal response body")
	}

	// Check for application-level errors.
	if r.Error != "" {
		return errors.Wrap(errors.New(r.Error), "API response contains error")
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

// Feed is a struct containing feed information from AdaFruit.
type Feed struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	LastValue   string `json:"last_value"`
	LastUpdated string `json:"last_value_at"`
	Key         string `json:"key"`
}
