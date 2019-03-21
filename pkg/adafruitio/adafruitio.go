// Package adafruitio provides an API client for Adafruit.IO.
package adafruitio

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type response struct {
	Error string `json:"error"`
}

// Do sets request headers, sends a request, checks for API errors, and
// returns the request body.
func Do(req *http.Request) ([]byte, error) {
	// Set headers.
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Send request.
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "could not send Adafruit API request")
	}
	defer res.Body.Close()

	// Parse response body.
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "could not read Adafruit API response body")
	}
	var r response
	err = json.Unmarshal(body, &r)
	if err != nil {
		// Special case: some API endpoints normally return arrays, but return
		// objects when an error occurs. In this case, unmarshalling an array result
		// to a struct (when the request succeeds) will fail. This is expected and
		// allowed behaviour.
		if err, ok := err.(*json.UnmarshalTypeError); !ok || err.Value != "array" {
			return nil, errors.Wrap(
				err, "could not unmarshal Adafruit API response body")
		}
	}

	// Check for application-level errors.
	if r.Error != "" {
		return nil, errors.Wrap(
			errors.New(r.Error), "Adafruit API response contains error")
	}

	return body, nil
}

// Authenticate adds an Adafruit API key header to an HTTP request.
func Authenticate(req *http.Request, key string) {
	req.Header.Set("X-AIO-Key", key)
}
