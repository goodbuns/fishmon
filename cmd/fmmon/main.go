// Fmmon is a fishmon monitor that checks Adafruit feed uptime and temperature values.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
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
	for {
		var feed []adafruitio.Feed
		feed, err = client.FeedsInGroup(*group)
		if err != nil {
			log.Fatalf("could not get feed information from Adafruit: %s", err.Error())
		}

		analyze(&feed, *numFeeds, *minTemp, *maxTemp)

		time.Sleep(time.Second * 3)
	}
}

func analyze(feeds *[]adafruitio.Feed, numFeeds int, minTemp, maxTemp float64) {
	// Check number of feeds.
	roleCall := false
	if len(*feeds) < numFeeds {
		alert("expected " + strconv.Itoa(numFeeds) + ", got " + strconv.Itoa(len(*feeds)))
		roleCall = true
	}

	// Check temperatures of feeds.
	for _, feed := range *feeds {
		if roleCall {
			alert(feed.Name + " still online")
		}
		temp, err := strconv.ParseFloat(feed.LastValue, 64)
		if err != nil {
			log.Println("could not parse feed value to float", err.Error())
		}
		if temp > maxTemp {
			alert(feed.LastUpdated + ": " + feed.Name + " temperature too high - " + feed.LastValue)
		}
		if temp < minTemp {
			alert(feed.LastUpdated + ": " + feed.Name + " temperature too low - " + feed.LastValue)
		}
	}
}

func alert(msg string) {
	// TODO: update with slack integration in alert PR
	log.Println(msg)
}
