// Copyright 2019 Adam Shannon
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/timpalpant/go-iex"
)

var (
	// Prometheus metrics
	stockAsks = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "stock_asks",
		Help: "..",
	}, []string{"symbol"})
	stockBids = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "stock_bids",
		Help: "..",
	}, []string{"symbol"})
	stockPrices = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "stock_prices",
		Help: "..",
	}, []string{"symbol"})
	stockVolumes = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "stock_volumes",
		Help: "..",
	}, []string{"symbol"})
	stockDataRefreshHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "stock_data_refresh_duration_milliseconds",
		Help:    "...",
		Buckets: []float64{10.0, 25.0, 50.0, 100.0, 250.0, 500.0, 1000.0, 2500.0, 5000.0},
	})
)

func init() {
	prometheus.MustRegister(stockAsks)
	prometheus.MustRegister(stockBids)
	prometheus.MustRegister(stockPrices)
	prometheus.MustRegister(stockVolumes)
	prometheus.MustRegister(stockDataRefreshHistogram)
}

func captureStockData(config *Config, iexClient *iex.Client) {
	symbols := config.Stocks.Symbols
	log.Printf("loading stock data for %s\n", strings.Join(symbols, ", "))
	for {
		start := time.Now()

		// TODO(adam): only run this when the market is open
		tops, err := iexClient.GetTOPS(symbols)
		if err != nil {
			log.Printf("ERROR: in TOPS: %v", err)
		} else {
			for i := range tops {
				t := tops[i]
				stockAsks.WithLabelValues(t.Symbol).Set(t.AskPrice)
				stockBids.WithLabelValues(t.Symbol).Set(t.BidPrice)
				stockVolumes.WithLabelValues(t.Symbol).Set(float64(t.Volume))
			}
		}

		quotes, err := iexClient.GetLast(symbols)
		if err != nil {
			panic(err)
		}
		for i := range symbols {
			stockPrices.WithLabelValues(symbols[i]).Set(quotes[i].Price)
		}

		// Capture loop time as milliseconds
		stockDataRefreshHistogram.Observe(float64(time.Since(start).Nanoseconds() / 1e6))

		time.Sleep(*flagInterval) // TODO(adam): When outside market hours turn this down to 5mins? or 30mins?
	}
}
