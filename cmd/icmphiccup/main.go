package main

import (
	"flag"
	"fmt"
	"github.com/fhalim/ghiccup"
	"github.com/tatsushid/go-fastping"
	"log"
	"net"
	"time"
)

func main() {
	interval := flag.Int("interval", 1, "Interval in seconds between pings")
	host := flag.String("host", "8.8.8.8", "Host to ping")
	rttThresholdMs := flag.Int("threshold", 100, "Threshold in ms above which to report RTTs")

	flag.Parse()

	threshold := time.Duration(*rttThresholdMs) * time.Millisecond

	p := initializePinger(*host, threshold)

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
func initializePinger(host string, reportingThreshold time.Duration) *fastping.Pinger {
	p := fastping.NewPinger()
	log.Println("Pinging host", host)
	ra, err := net.ResolveIPAddr("ip4:icmp", host)

	if err != nil {
		log.Panic("Unable to resolve host", err)
	}
	p.AddIPAddr(ra)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		log.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
		if rtt > reportingThreshold {
			fmt.Println(ghiccup.Marshall(ghiccup.HiccupInfo{
				Timestamp: time.Now().Format(time.RFC3339),
				Threshold: reportingThreshold.Nanoseconds(),
				Duration:  rtt.Nanoseconds(),
			}))
		}
	}
	p.OnIdle = func() {
		if p.Err() != nil {
			log.Println("Ping timed out")
		}
	}
	return p
}
