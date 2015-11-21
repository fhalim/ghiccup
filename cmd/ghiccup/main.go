/*
Measures hiccups in the platform and reports those that exceed the specified threshold
*/
package main

import (
	"flag"
	"fmt"
	"github.com/codahale/hdrhistogram"
	"github.com/fhalim/ghiccup"
	"log"
	"time"
	"encoding/json"
	"io/ioutil"
)

const DEFAULT_POLLING_FREQUENCY = 1000000
const DEFAULT_SNAPSHOT_FREQUENCY = 1e+10
const DEFAULT_PAUSE_THRESHOLD = DEFAULT_POLLING_FREQUENCY
const MAX_INTERESTING_LATENCY_MINUTES = 10

func main() {
	resolutionNs := flag.Int64("resolution", DEFAULT_POLLING_FREQUENCY, "Frequency in ns at which to test")
	thresholdNs := flag.Int64("threshold", DEFAULT_PAUSE_THRESHOLD, "Threshold in ns of pause to allow")
	snapshotFileName := flag.String("outfileName", "histogram_snapshot.json", "File in which to store histogram snapshots")
	snapshotFrequencyNs := flag.Int64("snapshotFrequency", DEFAULT_SNAPSHOT_FREQUENCY, "Frequency in ns to save snapshots")

	flag.Parse()

	duration := time.Duration(*resolutionNs) * time.Nanosecond
	threshold := time.Duration(*thresholdNs) * time.Nanosecond
	snapshotFrequency := time.Duration(*snapshotFrequencyNs) * time.Nanosecond

	hist := hdrhistogram.New(0, (time.Duration(MAX_INTERESTING_LATENCY_MINUTES) * time.Minute).Nanoseconds(), 5)
	go execute(duration, threshold, hist)
	saveSnapshots(hist, *snapshotFileName, snapshotFrequency)
	log.Println("Done!")
}

func execute(resolution time.Duration, threshold time.Duration, hist *hdrhistogram.Histogram) {
	allowedPause := resolution.Nanoseconds() + threshold.Nanoseconds()
	for {
		startTime := time.Now()
		time.Sleep(resolution)
		endTime := time.Now()
		durationNs := endTime.Sub(startTime).Nanoseconds()
		hiccupNs := durationNs - resolution.Nanoseconds()
		hist.RecordValue(hiccupNs)
		if durationNs > allowedPause {
			fmt.Println(ghiccup.Marshall(ghiccup.HiccupInfo{
				Timestamp:  time.Now().Format(time.RFC3339),
				Resolution: resolution.Nanoseconds(),
				Threshold:  threshold.Nanoseconds(),
				Duration:   hiccupNs,
			}))
		}
	}
}

func saveSnapshots(hist *hdrhistogram.Histogram, outputFileName string, frequency time.Duration){
	for {
		time.Sleep(frequency)
		snapshot := hist.Export()
		b, err := json.Marshal(*snapshot)
		if err != nil {
			log.Panic(err)
		}
		ioutil.WriteFile(outputFileName, b, 0644)
	}
}