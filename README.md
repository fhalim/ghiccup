# Introduction

Small utility to measure hiccups on the runtime. Inspired by [JHiccup](http://www.azulsystems.com/jhiccup).


# Requirements

- [Go](http://golang.org/)
- [InfluxDB](http://influxdb.com) (Optional)

# Building

Ensure that [`GOPATH`](https://golang.org/doc/code.html#GOPATH) is set up properly.


```bash
go get github.com/fhalim/ghiccup
go get github.com/fhalim/ghiccup/ghiccup2influxdb
```

# Usage

The `ghiccup` schedules pauses and measures latency above the specified duration. Commandline options allow specifying the pause duration as well as threshold for latency values that should be reported. The utility emits JSON messages, one per line, on its STDOUT.

The `ghiccup2influxdb` utility can be used to write the information emitted by `ghiccup` into an InfluxDB series.
