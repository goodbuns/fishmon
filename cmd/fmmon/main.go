// Fmmon is a fishmon monitor that checks Adafruit feed uptime and temperature values.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/goodbuns/fishmon/pkg/adafruitio"
	"github.com/goodbuns/fishmon/pkg/ds18b20"
)

func main() {
	// Set up command-line flags.
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), `%s starts the feed monitoring service..

Usage of %s:
`, os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}
	aioUser := flag.String("aio_username", "", "Adafruit.IO username")
	aioKey := flag.String("aio_key", "", "Adafruit.IO key")
	var minTemp, maxTemp float64
	flag.Float64Var(&minTemp, "min_temp", 65.0, "Lowest temperature allowed before alerting, in degrees Fahrenheit. Default is 65.0.")
	flag.Float64Var(&maxTemp, "max_temp", 83.0, "Highest temperature allowed before alerting, in degrees Fahrenheit. Default is 83.0.")

	// Set up Adafruit.IO client.
	client, err := adafruitio.New(*aioUser, *aioKey)
	if err != nil {
		log.Fatalf("could not set up Adafruit.IO client: %s", err.Error())
	}

	// Find all sensors.
	sensors, err := ds18b20.Sensors()
	if err != nil {
		log.Fatalf("could not detect DS18B20 sensors: %s", err.Error())
	}
	log.Printf("found %d sensors: %#v\n", len(sensors), sensors)

	// Monitor Adafruit feed uptime.
	for {
		var feed []adafruitio.Feed
		feed, err = client.GetFeed()
		if err != nil {
			log.Fatalf("could not get feed information from Adafruit: %s", err.Error())
		}
		analyze(&feed)
	}
}

func analyze(*[]adafruitio.Feed) {
	// update later, should analyze the feed length, temperatures etc. to determine if alerts need to be sent
}
