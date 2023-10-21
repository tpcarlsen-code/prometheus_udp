package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/tpcarlsen-code/prometheus_udp/bridge"
)

func main() {
	udpPort := env("UDP_PORT", 9230)
	httpPort := env("HTTP_PORT", 9231)

	metrics := bridge.NewMetrics()
	metricsServer := bridge.NewMetricsServer(metrics)
	nr, err := bridge.DefaultUDPNetworkReader(udpPort.(int))
	if err != nil {
		log.Fatal(err)
	}
	intakeServer := bridge.NewIntakeServer(metrics, nr)

	ec := make(chan error)
	// Start HTTP server.
	go func() {
		ec <- metricsServer.Run(httpPort.(int))
	}()

	// Start UDP server.
	go func() {
		ec <- intakeServer.Run()
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	select {
	case err = <-ec:
		log.Fatalf("received error: %s", err.Error())
	case sig := <-signals:
		log.Fatalf("received os signal: %s, exiting.", sig.String())
	}
}

func env(envName string, def any) any {
	arg := os.Getenv(envName)
	if arg != "" {
		switch def.(type) {
		case int:
			val, err := strconv.Atoi(arg)
			if err != nil {
				log.Fatal(err)
			}
			return val
		}
	}
	return def
}
