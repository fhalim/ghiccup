package main

import (
	"flag"
	"fmt"
	utils "github.com/fhalim/ghiccup/utils"
	"time"
)

func main() {
	resolutionNs := flag.Uint64("resolution", 1000000, "Frequency in ns at which to test")
	thresholdNs := flag.Uint64("threshold", 1000000, "Threshold in ns of pause to allow above resolution")

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
			fmt.Println(utils.Marshall(utils.HiccupInfo{
				Timestamp:  time.Now().Format(time.RFC3339),
				Resolution: resolution.Nanoseconds(),
				Threshold:  threshold.Nanoseconds(),
				Duration:   durationNs,
			}))
		}
	}
}