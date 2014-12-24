package main

import (
	"flag"
	"fmt"
	"github.com/fhalim/ghiccup"
	"time"
)

const DEFAULT_POLLING_FREQUENCY = 1000000
const DEFAULT_PAUSE_THRESHOLD = DEFAULT_POLLING_FREQUENCY

func main() {
	resolutionNs := flag.Int64("resolution", DEFAULT_POLLING_FREQUENCY, "Frequency in ns at which to test")
	thresholdNs := flag.Int64("threshold", DEFAULT_PAUSE_THRESHOLD, "Threshold in ns of pause to allow")

	flag.Parse()

	duration := time.Duration(*resolutionNs) * time.Nanosecond
	threshold := time.Duration(*resolutionNs+*thresholdNs) * time.Nanosecond

	execute(duration, threshold)
}

func execute(resolution time.Duration, threshold time.Duration) {
	for {
		startTime := time.Now()
		time.Sleep(resolution)
		endTime := time.Now()
		durationNs := endTime.Sub(startTime).Nanoseconds()
		if durationNs > threshold.Nanoseconds() {
			fmt.Println(ghiccup.Marshall(ghiccup.HiccupInfo{
				Timestamp:  time.Now().Format(time.RFC3339),
				Resolution: resolution.Nanoseconds(),
				Threshold:  threshold.Nanoseconds(),
				Duration:   durationNs,
			}))
		}
	}
}
