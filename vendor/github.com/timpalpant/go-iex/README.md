# go-iex
A Go library for accessing the IEX Developer API.

[![GoDoc](https://godoc.org/github.com/timpalpant/go-iex?status.svg)](http://godoc.org/github.com/timpalpant/go-iex)
[![Build Status](https://travis-ci.org/timpalpant/go-iex.svg?branch=master)](https://travis-ci.org/timpalpant/go-iex)
[![Coverage Status](https://coveralls.io/repos/timpalpant/go-iex/badge.svg?branch=master&service=github)](https://coveralls.io/github/timpalpant/go-iex?branch=master)

go-iex is a library to access the [IEX Developer API](https://www.iextrading.com/developer/docs/) from [Go](http://www.golang.org).
It provides a thin wrapper for working with the JSON REST endpoints and [IEXTP1 pcap data](https://www.iextrading.com/trading/market-data/#specifications).

[IEX](https://www.iextrading.com) is a fair, simple and transparent stock exchange dedicated to investor protection.
IEX provides realtime and historical market data for free through the IEX Developer API.
By using the IEX API, you agree to the [Terms of Use](https://www.iextrading.com/api-terms/). IEX is not affiliated
and does not endorse or recommend this library.

## Usage

### pcap2json

If you just need a tool to convert the provided pcap data files into JSON, you can use the included `pcap2json` tool:

```
$ go install github.com/timpalpant/go-iex/pcap2json
$ pcap2json < input.pcap > output.json
```

### Fetch real-time top-of-book quotes

```Go
package main

import (
  "fmt"
  "net/http"

  "github.com/timpalpant/go-iex"
)

func main() {
  client := iex.NewClient(&http.Client{})

  quotes, err := client.GetTOPS([]string{"AAPL", "SPY"})
  if err != nil {
      panic(err)
  }

  for _, quote := range quotes {
      fmt.Printf("%v: bid $%.02f (%v shares), ask $%.02f (%v shares) [as of %v]\n",
          quote.Symbol, quote.BidPrice, quote.BidSize,
          quote.AskPrice, quote.AskSize, quote.LastUpdated)
  }
}
```

### Fetch historical top-of-book quote (L1 tick) data.

Historical tick data (TOPS and DEEP) can be parsed using the `PcapScanner`.

```Go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/timpalpant/go-iex"
	"github.com/timpalpant/go-iex/iextp/tops"
)

func main() {
	client := iex.NewClient(&http.Client{})

	// Get historical data dumps available for 2016-12-12.
	histData, err := client.GetHIST(time.Date(2016, time.December, 12, 0, 0, 0, 0, time.UTC))
	if err != nil {
		panic(err)
	} else if len(histData) == 0 {
		panic(fmt.Errorf("Found %v available data feeds", len(histData)))
	}

	// Fetch the pcap dump for that date and iterate through its messages.
	resp, err := http.Get(histData[0].Link)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	packetDataSource, err := iex.NewPacketDataSource(resp.Body)
	if err != nil {
		panic(err)
	}
	pcapScanner := iex.NewPcapScanner(packetDataSource)

	// Write each quote update message to stdout, in JSON format.
	enc := json.NewEncoder(os.Stdout)

	for {
		msg, err := pcapScanner.NextMessage()
		if err != nil {
			if err == io.EOF {
				break
			}

			panic(err)
		}

		switch msg := msg.(type) {
		case *tops.QuoteUpdateMessage:
			enc.Encode(msg)
		default:
		}
	}
}
```

## Contributing

Pull requests and issues are welcomed!

## License

go-iex is released under the [GNU Lesser General Public License, Version 3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html)