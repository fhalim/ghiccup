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
)

const DEFAULT_POLLING_FREQUENCY = 1000000
const DEFAULT_PAUSE_THRESHOLD = DEFAULT_POLLING_FREQUENCY
const MAX_INTERESTING_LATENCY_MINUTES = 10

func main() {
	resolutionNs := flag.Int64("resolution", DEFAULT_POLLING_FREQUENCY, "Frequency in ns at which to test")
	thresholdNs := flag.Int64("threshold", DEFAULT_PAUSE_THRESHOLD, "Threshold in ns of pause to allow")

	flag.Parse()

	duration := time.Duration(*resolutionNs) * time.Nanosecond
	threshold := time.Duration(*resolutionNs+*thresholdNs) * time.Nanosecond

	hist := hdrhistogram.New(0, (time.Duration(MAX_INTERESTING_LATENCY_MINUTES) * time.Minute).Nanoseconds(), 5)
	execute(duration, threshold, hist)
}

func execute(resolution time.Duration, threshold time.Duration, hist *hdrhistogram.Histogram) {
	for {
		startTime := time.Now()
		time.Sleep(resolution)
		endTime := time.Now()
		durationNs := endTime.Sub(startTime).Nanoseconds()
		hist.RecordValue(durationNs)
		if durationNs > threshold.Nanoseconds() {
			fmt.Println(ghiccup.Marshall(ghiccup.HiccupInfo{
				Timestamp:  time.Now().Format(time.RFC3339),
				Resolution: resolution.Nanoseconds(),
				Threshold:  threshold.Nanoseconds(),
				Duration:   durationNs,
			}))
			log.Println("Percentiles: 95:", hist.ValueAtQuantile(95), "50:", hist.ValueAtQuantile(50))
		}
	}
}
