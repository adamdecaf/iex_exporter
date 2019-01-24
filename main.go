// Copyright 2019 Adam Shannon
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/timpalpant/go-iex"
	"gopkg.in/yaml.v2"
)

const Version = "0.1.2-dev"

var (
	defaultInterval, _ = time.ParseDuration("1m")

	// CLI flags
	flagAddress    = flag.String("address", "0.0.0.0:9099", "HTTP listen address")
	flagConfigFile = flag.String("config.file", "", "Path to config file")
	flagInterval   = flag.Duration("interval", defaultInterval, "Interval to check metrics at")
	flagVersion    = flag.Bool("version", false, "Print the iex_exporter version")
)

func main() {
	flag.Parse()

	if *flagVersion {
		fmt.Println(Version)
		os.Exit(1)
	}

	apiToken := os.Getenv("IEX_API_TOKEN")
	if apiToken == "" {
		panic("IEX_API_TOKEN is required!")
	}
	log.Printf("Starting iex_exporter (Version: %s)\n", Version)

	// Read our config file
	var config *Config
	if *flagConfigFile == "" {
		log.Println("-config.file is empty so using example config")
		config = &Config{
			Stocks: &StocksConfig{
				Symbols: []string{"AAPL", "FB"},
			},
		}
	} else {
		bs, err := ioutil.ReadFile(*flagConfigFile)
		if err != nil {
			log.Fatalf("problem reading %s: %v", *flagConfigFile, err)
		}
		if err := yaml.Unmarshal(bs, &config); err != nil {
			log.Fatalf("problem unmarshaling %s: %v", *flagConfigFile, err)
		}
	}

	// Bring up IEX client and exporters
	iexClient := iex.NewClient(&http.Client{
		Timeout: 5 * time.Second,
	})
	go captureStockData(config, iexClient)

	// Start Prometheus metrics endpoint
	h := promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{})
	http.Handle("/metrics", h)

	log.Printf("listenting on %s", *flagAddress)
	if err := http.ListenAndServe(*flagAddress, nil); err != nil {
		log.Fatalf("ERROR binding to %s: %v", *flagAddress, err)
	}
}
