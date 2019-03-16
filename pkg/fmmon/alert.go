package fmmon

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Alert sends the specified message to the specified webhook URL.
func Alert(msg string, webhookURL string) {
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
