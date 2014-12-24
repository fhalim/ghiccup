package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type HiccupInfo struct {
	Timestamp  string `json:"timestamp"`
	Resolution uint64 `json:"resolution"`
	Threshold  uint64 `json:"threshold"`
	Duration   uint64 `json:"duration"`
}

type InfluxDbMetricValue struct {
	Name    string        `json:"name"`
	Columns []string      `json:"columns"`
	Points  []interface{} `json:"points"`
}

func main() {
	bio := bufio.NewReader(os.Stdin)
	for {

		line, _, err := bio.ReadLine()

		hiccupInfo := readHiccupInfo(line)
		influxDbPayload := createInfluxDbPayload([]HiccupInfo{hiccupInfo})

		fmt.Printf(influxDbPayload + "\n")
		postData(influxDbPayload)

		if err != nil {
			break
		}
	}

}

func postData(payload string) {
	url := "http://localhost:8086/db/local_dev/series?u=test&p=test"
	resp, err := http.Post(url, "application/json", strings.NewReader(payload))
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		fmt.Println("Received error code", resp.StatusCode, resp)
	}
}

func createInfluxDbPayload(hiccupInfos []HiccupInfo) string {
	hostname, _ := os.Hostname()
	hiccupInfo := hiccupInfos[0]
	timestamp, _ := time.Parse(time.RFC3339, hiccupInfo.Timestamp)
	hiccupDuration := hiccupInfo.Duration - hiccupInfo.Resolution

	metricValue := InfluxDbMetricValue{
		Name:    "hiccups",
		Columns: []string{"time", "hiccupDuration", "hostname"},
		Points:  []interface{}{[]interface{}{timestamp.UnixNano() / 1000000, hiccupDuration, hostname}},
	}

	influxDbLine := marshall([]InfluxDbMetricValue{metricValue})
	return influxDbLine
}

func readHiccupInfo(line []byte) HiccupInfo {
	var hiccupInfo HiccupInfo
	err := json.Unmarshal(line, &hiccupInfo)

	if err != nil {
		panic(err)
	}

	return hiccupInfo
}
func marshall(value interface{}) string {
	bytes, err := json.Marshal(value)
	if err != nil {
		panic(fmt.Sprintf("Error %v", err))
	}
	return string(bytes)
}
