// Fmmon is a fishmon monitor that checks Adafruit feed uptime and temperature values.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/goodbuns/fishmon/pkg/adafruitio"
	"github.com/goodbuns/fishmon/pkg/fishmonbot"
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

	// Set up Adafruit.IO client.
	client, err := adafruitio.New(*aioUser, "")
	if err != nil {
		log.Fatalf("could not set up Adafruit.IO client: %s", err.Error())
	}

	// Monitor Adafruit feed uptime.
	var feed []adafruitio.Feed
	feed, err = client.FeedsInGroup(*group)
	alerts := make(chan string, 6)

	// Sends updates every 5 hours.
	go func() {
		time.Sleep(time.Second * 2)
		for {
			fishmonbot.Analyze(&feed, alerts, *numFeeds, *minTemp, *maxTemp, *webhookURL, true)
			time.Sleep(time.Hour * 5)
		}
	}()

	// Rate limits alerts to every 5 hours.
	go func(alerts chan string, webhookURL string, start bool) {
		for {
			msg := ""
			select {
			case m := <-alerts:
				msg = "------------------------" + "\n" + m
				for i := 0; i < 3; i++ {
					msg = msg + "\n" + <-alerts
				}
				msg = msg + "\n------------------------"
				fishmonbot.Alert(msg, webhookURL)
				start = false
			default:
				if start {
					continue
				}
			}
			log.Println("SLEEPING NOW")
			time.Sleep(time.Hour * 5)
			start = true
		}
	}(alerts, *webhookURL, true)

	// Checks status every minute.
	for {
		feed, err = client.FeedsInGroup(*group)
		if err != nil {
			fishmonbot.Alert(":alarm: Could not get feed information from Adafruit - feeds may be down!", *webhookURL)
			log.Fatalf("could not get feed information from Adafruit: %s", err.Error())
		}

		fishmonbot.Analyze(&feed, alerts, *numFeeds, *minTemp, *maxTemp, *webhookURL, false)
		time.Sleep(time.Second)
	}
}
