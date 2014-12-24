package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	ghiccup "github.com/fhalim/ghiccup/utils"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type InfluxDbMetricValue struct {
	Name    string        `json:"name"`
	Columns []string      `json:"columns"`
	Points  []interface{} `json:"points"`
}

func main() {
	baseUrl := flag.String("baseurl", "http://localhost:8086", "Base URL for InfluxDB server")
	username := flag.String("username", "admin", "InfluxDB username")
	password := flag.String("password", "admin", "InfluxDB password")
	database := flag.String("database", "local_dev", "InfluxDB database")
	series := flag.String("series", "hiccups", "Data series")

	flag.Parse()

	url := fmt.Sprintf("%v/db/%v/series?u=%v&p=%v", *baseUrl, *database, *username, *password)

	bio := bufio.NewReader(os.Stdin)
	for {

		line, _, err := bio.ReadLine()

		hiccupInfo := readHiccupInfo(line)
		influxDbPayload := createInfluxDbPayload([]ghiccup.HiccupInfo{hiccupInfo}, *series)

		log.Println("Sending payload", influxDbPayload)
		postData(influxDbPayload, url)

		if err != nil {
			break
		}
	}

}

func postData(payload string, url string) {
	const MIME_TYPE = "application/json"
	resp, err := http.Post(url, MIME_TYPE, strings.NewReader(payload))
	if err != nil {
		log.Panicln("Error from HTTP request", err)
	}
	if resp.StatusCode != 200 {
		log.Panicln("Received error code", resp.StatusCode, resp)
	}
}

func createInfluxDbPayload(hiccupInfos []ghiccup.HiccupInfo, series string) string {
	hostname, _ := os.Hostname()
	hiccupInfo := hiccupInfos[0]
	timestamp, _ := time.Parse(time.RFC3339, hiccupInfo.Timestamp)
	hiccupDuration := hiccupInfo.Duration - hiccupInfo.Resolution

	metricValue := InfluxDbMetricValue{
		Name:    series,
		Columns: []string{"time", "hiccupDuration", "hostname"},
		Points:  []interface{}{[]interface{}{timestamp.UnixNano() / 1000000, hiccupDuration, hostname}},
	}

	influxDbLine := ghiccup.Marshall([]InfluxDbMetricValue{metricValue})
	return influxDbLine
}

func readHiccupInfo(line []byte) ghiccup.HiccupInfo {
	var hiccupInfo ghiccup.HiccupInfo
	err := json.Unmarshal(line, &hiccupInfo)

	if err != nil {
		log.Panicln("Error reading HiccupInfo", err)
	}

	return hiccupInfo
}
