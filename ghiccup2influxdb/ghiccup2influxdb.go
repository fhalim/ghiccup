package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	ghiccup "github.com/fhalim/ghiccup/utils"
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
	bio := bufio.NewReader(os.Stdin)
	for {

		line, _, err := bio.ReadLine()

		hiccupInfo := readHiccupInfo(line)
		influxDbPayload := createInfluxDbPayload([]ghiccup.HiccupInfo{hiccupInfo})

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

func createInfluxDbPayload(hiccupInfos []ghiccup.HiccupInfo) string {
	hostname, _ := os.Hostname()
	hiccupInfo := hiccupInfos[0]
	timestamp, _ := time.Parse(time.RFC3339, hiccupInfo.Timestamp)
	hiccupDuration := hiccupInfo.Duration - hiccupInfo.Resolution

	metricValue := InfluxDbMetricValue{
		Name:    "hiccups",
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
		panic(err)
	}

	return hiccupInfo
}
