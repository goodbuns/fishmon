// Fmmon is a fishmon monitor that checks Adafruit feed uptime and temperature values.
package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/goodbuns/fishmon/pkg/adafruitio"
)

type Sample struct {
	ActualNumFeeds   int
	ExpectedNumFeeds int

	CouldNotRetrieveData map[adafruitio.FeedID]bool
	CouldNotParseData    map[adafruitio.FeedID]bool
	BelowMinTemp         map[adafruitio.FeedID]bool
	AboveMaxTemp         map[adafruitio.FeedID]bool
	Stale                map[adafruitio.FeedID]bool

	Feeds map[adafruitio.FeedID]adafruitio.Feed
}

func (s *Sample) String() string {
	var alarms []string

	if s.ActualNumFeeds != s.ExpectedNumFeeds {
		alarms = append(alarms, fmt.Sprintf(":alarm: Expected %d feeds, but found %d instead", s.ExpectedNumFeeds, s.ActualNumFeeds))
	}
	for id := range s.CouldNotRetrieveData {
		feed := s.Feeds[id]
		alarms = append(alarms, fmt.Sprintf(":alarm: Could not retrieve data for feed %s (%s)", feed.Key, feed.Name))
	}
	for id := range s.CouldNotParseData {
		feed := s.Feeds[id]
		alarms = append(alarms, fmt.Sprintf(":alarm: Could not parse data for feed %s (%s)", feed.Key, feed.Name))
	}
	for id := range s.BelowMinTemp {
		feed := s.Feeds[id]
		alarms = append(alarms, fmt.Sprintf(":alarm: :thermometer: %s is below minimum temperature", feed.Name))
	}
	for id := range s.AboveMaxTemp {
		feed := s.Feeds[id]
		alarms = append(alarms, fmt.Sprintf(":alarm: :thermometer: %s is above maximum temperature", feed.Name))
	}
	for id := range s.Stale {
		feed := s.Feeds[id]
		alarms = append(alarms, fmt.Sprintf(":alarm: %s probes are not reporting", feed.Name))
	}

	var temperatures []string
	for _, feed := range s.Feeds {
		temperatures = append(temperatures, fmt.Sprintf(":thermometer: %s - %s°F", feed.Name, feed.LastValue))
	}
	sort.Slice(temperatures, func(i, j int) bool {
		return temperatures[i] < temperatures[j]
	})

	message := strings.Join(append(alarms, temperatures...), "\n")
	if len(alarms) == 0 {
		message = ":heavy_check_mark: OK\n" + message
	} else {
		message = ":alarm: <!channel>\n" + message
	}

	return message
}

func main() {
	// Set up command-line flags.
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), `%s starts the fmmon service.

Fmmon requests feed information from Adafruit to monitor and alert on feed uptime and tank temperature values.
	
Usage of %s:
`, os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}
	user := flag.String("user", "", "Adafruit.IO username")
	group := flag.String("group", "fish", "Name of Adafruit.IO group feeds to monitor")
	expectedNumFeeds := flag.Int("expected_num_feeds", 0, "Expected number of online feeds within the specified group")
	minTemp := flag.Float64("min_temp", 65, "Lowest temperature allowed before alerting, in degrees Fahrenheit")
	maxTemp := flag.Float64("max_temp", 83, "Highest temperature allowed before alerting, in degrees Fahrenheit")
	pollInterval := flag.Int("poll", 5*60, "Polling interval, in seconds")
	webhookURL := flag.String("webhook_url", "", "Webhook URL")
	flag.Parse()

	// Monitor Adafruit feed uptime.
	for {
		feeds, err := adafruitio.Group(*user, *group)
		if err != nil {
			Send(*webhookURL, fmt.Sprintf("Could not get feed group: %s", err.Error()))
			time.Sleep(5 * time.Minute)
			continue
		}

		sample := Sample{
			ExpectedNumFeeds:     *expectedNumFeeds,
			ActualNumFeeds:       len(feeds),
			CouldNotRetrieveData: make(map[adafruitio.FeedID]bool),
			CouldNotParseData:    make(map[adafruitio.FeedID]bool),
			BelowMinTemp:         make(map[adafruitio.FeedID]bool),
			AboveMaxTemp:         make(map[adafruitio.FeedID]bool),
			Stale:                make(map[adafruitio.FeedID]bool),
			Feeds:                make(map[adafruitio.FeedID]adafruitio.Feed),
		}

		last := time.Now().Add(-10 * time.Minute)
		for _, feed := range feeds {
			sample.Feeds[feed.ID] = feed

			// Check for liveness.
			if feed.LastUpdated.After(last) {
				sample.Stale[feed.ID] = true
			}

			// Retrieve temperature readings.
			points, err := adafruitio.Data(*user, feed.Key, last)
			if err != nil {
				sample.CouldNotRetrieveData[feed.ID] = true
				continue
			}

			// Check temperature readings.
			for _, point := range points {
				value, err := strconv.ParseFloat(point.Value, 64)
				if err != nil {
					sample.CouldNotParseData[feed.ID] = true
				}
				if value < *minTemp {
					sample.BelowMinTemp[feed.ID] = true
				}
				if value > *maxTemp {
					sample.AboveMaxTemp[feed.ID] = true
				}
			}
		}

		Send(*webhookURL, sample.String())
		time.Sleep(time.Duration(*pollInterval) * time.Second)
	}
}
