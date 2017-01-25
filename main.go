package main

func main() {
	ticker := Ticker{
		BtcUsd:              make(chan float64, 10),
		BtcEur:              make(chan float64, 10),
		UsdActiveFeedsCount: make(chan int, 10),
		EurActiveFeedsCount: make(chan int, 10),
	}

	ticker.SetFeeds([]BitcoinFeed{
		// emulating three different BTC-USD feeds
		&Gdax{Currency: "BTC-USD", Output: ticker.BtcUsd, FeedsCount: ticker.UsdActiveFeedsCount},
		&Gdax{Currency: "BTC-USD", Output: ticker.BtcUsd, FeedsCount: ticker.UsdActiveFeedsCount},
		&Gdax{Currency: "BTC-USD", Output: ticker.BtcUsd, FeedsCount: ticker.UsdActiveFeedsCount},
		// emulating four different BTC-EUR feeds
		&Gdax{Currency: "BTC-EUR", Output: ticker.BtcEur, FeedsCount: ticker.EurActiveFeedsCount},
		&Gdax{Currency: "BTC-EUR", Output: ticker.BtcEur, FeedsCount: ticker.EurActiveFeedsCount},
		&Gdax{Currency: "BTC-EUR", Output: ticker.BtcEur, FeedsCount: ticker.EurActiveFeedsCount},
		&Gdax{Currency: "BTC-EUR", Output: ticker.BtcEur, FeedsCount: ticker.EurActiveFeedsCount},
	})

	for _, feed := range ticker.Feeds {
		go feed.Subscribe()
	}

	ticker.Print()
}
