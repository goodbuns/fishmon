package adafruitio

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// A FeedID uniquely identifies an Adafruit feed.
type FeedID int

// Feed contains Adafruit feed metadata.
type Feed struct {
	ID          FeedID    `json:"id"`
	Name        string    `json:"name"`
	Key         string    `json:"key"`
	LastValue   string    `json:"last_value"`
	LastUpdated time.Time `json:"last_value_at"`
}

// Feeds retrieves all public feeds of an Adafruit user.
func Feeds(user string) ([]Feed, error) {
	// Construct request.
	req, err := http.NewRequest(
		http.MethodGet,
		"https://io.adafruit.com/api/v2/"+user+"/feeds",
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "could not construct feeds API request")
	}

	// Send request.
	res, err := Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "could not send feeds API request")
	}

	// Unmarshal response into Feed.
	var feeds []Feed
	err = json.Unmarshal(res, &feeds)
	if err != nil {
		return nil, errors.Wrap(
			err, "could not unmarshal response body for feed data")
	}

	return feeds, nil
}

// Group retrieves all feeds in a group of an Adafruit user.
func Group(user, group string) ([]Feed, error) {
	// Construct request.
	req, err := http.NewRequest(
		http.MethodGet,
		"https://io.adafruit.com/api/v2/"+user+"/groups/"+group+"/feeds",
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "could not construct group feed API request")
	}

	// Send request.
	res, err := Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "could not send group feed API request")
	}

	// Unmarshal response into Feed.
	var feeds []Feed
	err = json.Unmarshal(res, &feeds)
	if err != nil {
		return nil, errors.Wrap(
			err, "could not unmarshal response body for group feed data")
	}

	return feeds, nil
}

// A Point is a single value from an Adafruit.IO feed.
type Point struct {
	ID        string    `json:"id"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Data retrieves all points in a feed since a moment in time.
func Data(user, feed string, since time.Time) ([]Point, error) {
	// Construct request.
	req, err := http.NewRequest(
		http.MethodGet,
		"https://io.adafruit.com/api/v2/"+user+"/feeds/"+feed+"/data",
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "could not construct data point API request")
	}

	// Send request.
	res, err := Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "could not send data point API request")
	}

	// Unmarshal response into Point.
	var points []Point
	err = json.Unmarshal(res, &points)
	if err != nil {
		return nil, errors.Wrap(
			err, "could not unmarshal response body for data point")
	}

	return points, nil
}
