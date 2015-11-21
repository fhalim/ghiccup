[![Build Status](https://travis-ci.org/fhalim/ghiccup.svg?branch=master)](https://travis-ci.org/fhalim/ghiccup)

# Introduction

Small utility to measure hiccups on the runtime. Inspired by [JHiccup](http://www.azulsystems.com/jhiccup).


# Requirements

- [Go](http://golang.org/)
- [InfluxDB](http://influxdb.com) (Optional)

# Building

Ensure that [`GOPATH`](https://golang.org/doc/code.html#GOPATH) is set up properly.


```bash
go get github.com/fhalim/ghiccup/cmd/ghiccup
go get github.com/fhalim/ghiccup/cmd/ghiccup2influxdb
go get github.com/fhalim/ghiccup/cmd/icmphiccup
```

# Usage

The `ghiccup` command schedules pauses and measures latency above the specified duration. Commandline options allow specifying the pause duration as well as threshold for latency values that should be reported. The utility emits JSON messages, one per line, on its STDOUT.

`icmphiccup` pings a specified host and reports in the same JSON format as `ghiccup` when the RTT exceeds a specified threshold.

The `ghiccup2influxdb` utility can be used to write the information emitted by `ghiccup` into an InfluxDB series.
