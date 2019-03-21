package adafruitio

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// A Client provides authentication for authenticated Adafruit.IO API requests.
type Client struct {
	username string
	apiKey   string
}

// New constructs an authenticated Adafruit.IO API client, checking to make sure
// that the credentials are valid.
func New(username, apiKey string) (*Client, error) {
	// Construct client.
	client := &Client{
		username: username,
		apiKey:   apiKey,
	}

	// Construct API request.
	req, err := http.NewRequest(
		http.MethodGet, "https://io.adafruit.com/api/v2/user", nil)
	if err != nil {
		return nil, errors.Wrap(
			err, "could not construct API request to validate credentials")
	}

	// Check client credentials.
	Authenticate(req, apiKey)
	_, err = Do(req)
	if err != nil {
		return nil, errors.Wrap(
			err, "API response for validating credentials has error")
	}

	return client, nil
}

// A DataRequest contains an Adafruit feed data point.
type DataRequest struct {
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
}

// Record uploads a value to an Adafruit.IO feed.
func (c *Client) Record(feed, value string, timestamp time.Time) error {
	// Marshal request body.
	payload, err := json.Marshal(DataRequest{
		Value:     value,
		CreatedAt: timestamp,
	})
	if err != nil {
		return errors.Wrap(
			err, "could not marshal API request body for recording data")
	}

	// Construct request.
	req, err := http.NewRequest(
		http.MethodPost,
		"https://io.adafruit.com/api/v2/"+c.username+"/feeds/"+feed+"/data",
		bytes.NewReader(payload),
	)
	if err != nil {
		return errors.Wrap(
			err, "could not construct API request for recording data")
	}

	// Send request.
	Authenticate(req, c.apiKey)
	_, err = Do(req)
	if err != nil {
		return errors.Wrap(err, "API response for recording data has error")
	}

	return nil
}
