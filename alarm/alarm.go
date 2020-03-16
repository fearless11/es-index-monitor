package alarm

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const contentTypeJSON = "application/json"

// Alert ...
type Alert struct {
	url       string
	Alertname string `json:"alertname"`
	From      string `json:"from"`
	Level     string `json:"level"`
	Txt       string `json:"txt"`
}

// New ...
func New(url, from, level string) *Alert {
	return &Alert{
		url:   url,
		From:  from,
		Level: level,
	}
}

// Send ...
func (a *Alert) Send(alertname string, txt string) {
	a.Alertname = alertname
	a.Txt = txt

	log.Println("WARN Alert", alertname, txt)

	tr := &http.Transport{
		MaxIdleConns:       30,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(*a); err != nil {
		log.Println("ERROR ", err)
	}

	c := &http.Client{Transport: tr}
	_, err := c.Post(a.url, contentTypeJSON, &buf)
	if err != nil {
		log.Println("ERROR ", err)
	}
}
