package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/timpalpant/go-iex"
)

const Version = "0.1.0-dev"

var (
	defaultInterval, _ = time.ParseDuration("1m")

	// CLI flags
	flagAddress = flag.String("address", "0.0.0.0:9099", "HTTP listen address")
	// flagConfigFile
	flagInterval = flag.Duration("interval", defaultInterval, "Interval to check domains at")
	flagVersion  = flag.Bool("version", false, "Print the rdap_exporter version")

	// Prometheus metrics
	stockQuotes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "stock_quotes",
			Help: "..",
		},
		[]string{"price"},
	)
)

func init() {
	prometheus.MustRegister(stockQuotes)
}

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

	// Bring up IEX client
	iexClient := iex.NewClient(&http.Client{
		Timeout: 5 * time.Second,
	})

	symbols := []string{"AAPL", "FB"}
	go func() {
		for {
			quotes, err := iexClient.GetLast(symbols)
			if err != nil {
				panic(err)
			}
			for i := range symbols {
				stockQuotes.WithLabelValues(symbols[i]).Set(quotes[i].Price)
			}
			time.Sleep(*flagInterval)
		}
	}()

	// Start Prometheus metrics endpoint
	h := promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{})
	http.Handle("/metrics", h)

	log.Printf("listenting on %s", *flagAddress)
	if err := http.ListenAndServe(*flagAddress, nil); err != nil {
		log.Fatalf("ERROR binding to %s: %v", *flagAddress, err)
	}
}
