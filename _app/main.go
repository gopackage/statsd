package main

import "github.com/gopackage/statsd"

func main() {
	server := statsd.Server{}
	server.Start()
}
