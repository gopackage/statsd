package main

import (
	"flag"

	"github.com/apex/log"
	"github.com/gopackage/statsd/server"
)

var udpAddr = flag.String("udp-addr", ":9045", "address to bind the UDP stats server")
var tcpAddr = flag.String("tcp-addr", "", "address to bind the TCP stats server if desired")
var httpAddr = flag.String("http-addr", ":8080", "address to bind the web server to view stats")
var respAddr = flag.String("resp-addr", "", "address to bind the RESP (Redis) server if desired")

func main() {
	flag.Parse()
	engine := server.New()
	log.WithFields(log.Fields{"udp":*udpAddr, "tcp":*tcpAddr, "http": *httpAddr, "resp":*respAddr}).Info("starting engine")
	err := engine.Start(*udpAddr, *tcpAddr, *httpAddr, *respAddr)
	if err != nil {
		log.WithError(err).Fatal("engine error")
	}
	log.Info("engine stopped successfully")
}
