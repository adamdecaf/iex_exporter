// Copyright 2019 Adam Shannon
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"strings"
	"sync"
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
		Buckets: []float64{10.0, 25.0, 50.0, 100.0, 250.0, 500.0, 750.0, 1000.0, 2500.0, 5000.0},
	})
)

func init() {
	prometheus.MustRegister(stockAsks)
	prometheus.MustRegister(stockBids)
	prometheus.MustRegister(stockPrices)
	prometheus.MustRegister(stockVolumes)
	prometheus.MustRegister(stockDataRefreshHistogram)

	loc, _ := time.LoadLocation("America/New_York")
	if estLocation == nil {
		estLocation = loc
	}
}

var estLocation *time.Location

func marketOpen(now time.Time) bool {
	now = now.In(estLocation)
	if now.Weekday() == time.Sunday || now.Weekday() == time.Saturday {
		return false
	}
	if (now.Hour() < 9 || now.Hour() > 16) || (now.Hour() == 16 && now.Minute() > 30) { // 9am to 4:30pm
		return false
	}
	return true
}

func captureStockData(config *Config, iexClient *iex.Client) {
	symbols := config.Stocks.Symbols
	log.Printf("loading stock data for %s\n", strings.Join(symbols, ", "))
	for {
		if !marketOpen(time.Now()) {
			continue
		}

		start := time.Now()
		wg := sync.WaitGroup{}

		// Capture stock metadata
		wg.Add(1)
		go func() {
			defer wg.Done()
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
		}()

		// Get stock quote
		wg.Add(1)
		go func() {
			defer wg.Done()
			quotes, err := iexClient.GetLast(symbols)
			if err != nil {
				log.Printf("ERROR: in GetLast: %v", err)
			} else {
				for i := range symbols {
					stockPrices.WithLabelValues(symbols[i]).Set(quotes[i].Price)
				}
			}
		}()

		wg.Wait()
		diff := time.Since(start).Nanoseconds()
		stockDataRefreshHistogram.Observe(float64(diff / 1e6))
		time.Sleep(*flagInterval - time.Duration(diff))
	}
}
