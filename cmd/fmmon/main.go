// Fmmon is a fishmon monitor that checks Adafruit feed uptime and temperature values.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/goodbuns/fishmon/pkg/adafruitio"
)

func main() {
	// Set up command-line flags.
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), `%s starts the fmmon service.

Fmmon requests feed information from Adafruit to monitor and alert on feed uptime and tank temperature values.
	
Usage of %s:
`, os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}
	aioUser := flag.String("aio_username", "", "Adafruit.IO username")
	// aioKey := flag.String("aio_key", "", "Adafruit.IO key")
	group := flag.String("group", "fish", "Group name of feeds to monitor.")
	numFeeds := flag.Int("num_feeds", 0, "Expected number of online feeds within the specified group.")
	minTemp := flag.Float64("min_temp", 65, "Lowest temperature allowed before alerting, in degrees Fahrenheit.")
	maxTemp := flag.Float64("max_temp", 83, "Highest temperature allowed before alerting, in degrees Fahrenheit.")
	webhookURL := flag.String("webhook_url", "", "Webhook URL.")
	flag.Parse()

	log.Println("min temp set", *minTemp)
	log.Println("max temp set", *maxTemp)
	log.Println(*aioUser)
	// log.Println(*aioKey)

	// Set up Adafruit.IO client.
	client, err := adafruitio.New(*aioUser, "")
	if err != nil {
		log.Fatalf("could not set up Adafruit.IO client: %s", err.Error())
	}

	// Monitor Adafruit feed uptime.
	var feed []adafruitio.Feed
	feed, err = client.FeedsInGroup(*group)

	go func() {
		for {
			analyze(&feed, *numFeeds, *minTemp, *maxTemp, *webhookURL, true)
			time.Sleep(time.Hour * 5)
		}
	}()

	for {
		feed, err = client.FeedsInGroup(*group)
		if err != nil {
			alert(":alarm: Could not get feed information from Adafruit - feeds may be down!", *webhookURL)
			log.Fatalf("could not get feed information from Adafruit: %s", err.Error())
		}

		analyze(&feed, *numFeeds, *minTemp, *maxTemp, *webhookURL, false)

		time.Sleep(time.Second * 3)
	}
}

func analyze(feeds *[]adafruitio.Feed, numFeeds int, minTemp, maxTemp float64, webhookURL string, update bool) {
	// Check number of feeds.
	roleCall := false
	if len(*feeds) < numFeeds {
		alert("expected "+strconv.Itoa(numFeeds)+"feeds, got only "+strconv.Itoa(len(*feeds)), webhookURL)
		roleCall = true
	}

	// Check temperatures of feeds.
	for _, feed := range *feeds {
		if roleCall {
			alert(feed.Name+" still online", webhookURL)
		}
		temp, err := strconv.ParseFloat(feed.LastValue, 64)
		if err != nil {
			log.Println("could not parse feed value to float", err.Error())
		}
		if temp > maxTemp {
			alert(":alarm:     "+feed.LastUpdated+": "+feed.Name+" :thermometer: temperature too high - "+feed.LastValue+"F", webhookURL)
		}
		if temp < minTemp {
			alert(":alarm:     "+feed.LastUpdated+": "+feed.Name+" :thermometer: temperature too low - "+feed.LastValue+"F", webhookURL)
		}
		if update {
			alert(feed.Name+" :thermometer: temperature - "+feed.LastValue, webhookURL)
		}
	}
}

func alert(msg string, webhookURL string) {
	m := `{"text":"` + msg + `"}`

	// Create webhook request.
	req, err := http.NewRequest("POST", webhookURL, strings.NewReader(m))
	req.Header.Set("Content-type", "application/json")
	if err != nil {
		log.Println("could not create HTTP POST request to webhook", err.Error())
	}
	log.Println(req)

	// Send webhook request.
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("could not post to webhook URL", err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("could not read response body", err.Error())
	}

	log.Println(string(body))

	defer res.Body.Close()
}
