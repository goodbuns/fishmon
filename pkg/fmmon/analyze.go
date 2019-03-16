package fmmon

import (
	"log"
	"strconv"

	"github.com/goodbuns/fishmon/pkg/adafruitio"
)

// Analyze analyzes feed information and alerts on any issues.
func Analyze(feeds *[]adafruitio.Feed, alerts chan string, numFeeds int, minTemp, maxTemp float64, webhookURL string, update bool) {
	// Check number of feeds.
	roleCall := false
	if len(*feeds) < numFeeds {
		m := "expected " + strconv.Itoa(numFeeds) + " feeds, got only " + strconv.Itoa(len(*feeds))
		alerts <- m
		roleCall = true
	}

	// Check temperatures of feeds.
	u := ""
	for _, feed := range *feeds {
		if roleCall {
			m := feed.Name + " still online"
			alerts <- m
		}
		temp, err := strconv.ParseFloat(feed.LastValue, 64)
		if err != nil {
			log.Println("could not parse feed value to float", err.Error())
		}
		if temp > maxTemp {
			m := ":alarm:     " + feed.LastUpdated + ": " + feed.Name + " :thermometer: temperature too high - " + feed.LastValue + "F"
			alerts <- m
		}
		if temp < minTemp {
			m := ":alarm:     " + feed.LastUpdated + ": " + feed.Name + " :thermometer: temperature too low - " + feed.LastValue + "F"
			alerts <- m
		}
		if update {
			u = u + "\n" + feed.Name + " :thermometer: temperature - " + feed.LastValue
		}
	}
	if update {
		Alert(u, webhookURL)
	}
}
