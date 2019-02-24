// Fishmon is a fish tank monitoring system for the Raspberry Pi.
package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/liftM/fishmon/config"
	"github.com/liftM/fishmon/pkg/adafruitio"
	"github.com/liftM/fishmon/pkg/ds18b20"
)

// Configurable constants.
const (
	RateLimitPerMinute = 30
)

func main() {
	rand.Seed(time.Now().Unix())

	// Set up command-line flags.
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), `%s starts the fishmon service.

Fishmon reads the outputs of connected DS18B20 temperature probes and uploads
them to Adafruit.IO.

Usage of %s:
`, os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}
	aioUser := flag.String("aio_username", "", "Adafruit.IO username")
	aioKey := flag.String("aio_key", "", "Adafruit.IO key")
	configFile := flag.String("config", "fishmonconfig.json", "Fishmon configuration file")
	flag.Parse()

	// Parse configuration.
	conf, err := config.New(*configFile)
	if err != nil {
		log.Fatalf("could not parse configuration file at %s: %s", *configFile, err.Error())
	}

	// Set up system.
	err = ds18b20.Ensure()
	if err != nil {
		log.Fatalf("could not set up DS18B20 probe: %s", err.Error())
	}

	// Set up sensors.
	sensors, err := ds18b20.Sensors()
	if err != nil {
		log.Fatalf("could not detect DS18B20 sensors: %s", err.Error())
	}
	log.Printf("found %d sensors: %#v\n", len(sensors), sensors)

	var probes []*ds18b20.Probe
	for _, sensor := range sensors {
		probe, err := ds18b20.New(sensor)
		if err != nil {
			log.Fatalf("could not set up probe %s: %s", sensor, err.Error())
		}
		probes = append(probes, probe)
	}

	// Set up Adafruit.IO client.
	client, err := adafruitio.New(*aioUser, *aioKey)
	if err != nil {
		log.Fatalf("could not set up Adafruit.IO client: %s", err.Error())
	}

	// Monitor and report temperature data.
	// Adafruit.IO limits free accounts to 30 data points per minute.
	rate := time.Minute / (RateLimitPerMinute / time.Duration(len(probes)))
	ticker := time.NewTicker(rate)

	for range ticker.C {
		for _, probe := range probes {
			timestamp := time.Now()

			// Sense temperature.
			temperature, err := probe.Sense()
			if err != nil {
				log.Printf("failed to sense temperature for probe %s: %s\n", probe.ID, err.Error())
				break
			}

			// Report temperature.
			pconf, ok := conf.Probes[probe.ID]
			if !ok {
				log.Fatalf("could not find configuration for probe %s", probe.ID)
			}
			client.Record(pconf.FeedKey, fmt.Sprintf("%.3f", temperature.Fahrenheit()), timestamp)

			fmt.Printf("time=%s probe=%s temp=%0.3f\n", timestamp.String(), probe.ID, temperature.Fahrenheit())
		}
	}
}
