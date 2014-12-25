package main

import (
	"flag"
	"fmt"
	"github.com/codahale/hdrhistogram"
	"github.com/fhalim/ghiccup"
	"github.com/tatsushid/go-fastping"
	"log"
	"net"
	"time"
)

const MAX_INTERESTING_LATENCY_SECONDS = 60

func main() {
	interval := flag.Int("interval", 1, "Interval in seconds between pings")
	host := flag.String("host", "8.8.8.8", "Host to ping")
	rttThresholdMs := flag.Int("threshold", 100, "Threshold in ms above which to report RTTs")

	flag.Parse()

	threshold := time.Duration(*rttThresholdMs) * time.Millisecond

	hist := hdrhistogram.New(0, (time.Duration(MAX_INTERESTING_LATENCY_SECONDS) * time.Second).Nanoseconds(), 5)

	p := initializePinger(*host, threshold, hist)

	go startPinger(p, time.Second*time.Duration(*interval))

	var input string
	fmt.Scanln(&input)
}

func startPinger(pinger *fastping.Pinger, frequency time.Duration) {
	for {
		go func() {
			err := pinger.Run()
			if err != nil {
				log.Panicln("Unable to send ping", err)
			}
		}()
		time.Sleep(frequency)
	}
}
func initializePinger(host string, reportingThreshold time.Duration, hist *hdrhistogram.Histogram) *fastping.Pinger {
	p := fastping.NewPinger()
	log.Println("Pinging host", host)
	ra, err := net.ResolveIPAddr("ip4:icmp", host)

	if err != nil {
		log.Panic("Unable to resolve host", err)
	}
	p.AddIPAddr(ra)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		log.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
		hist.RecordValue(rtt.Nanoseconds())
		if rtt > reportingThreshold {
			fmt.Println(ghiccup.Marshall(ghiccup.HiccupInfo{
				Timestamp: time.Now().Format(time.RFC3339),
				Threshold: reportingThreshold.Nanoseconds(),
				Duration:  rtt.Nanoseconds(),
			}))
			log.Println("Percentiles: 95:", hist.ValueAtQuantile(95), "50:", hist.ValueAtQuantile(50))
		}
	}
	p.OnIdle = func() {
		if p.Err() != nil {
			log.Println("Ping timed out")
		}
	}
	return p
}
