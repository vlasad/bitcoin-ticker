package main

const btcusd = "BTC-USD"
const btceur = "BTC-EUR"

func main() {
	ticker := Ticker{
		BtcUsd:              make(chan float64, 10),
		BtcEur:              make(chan float64, 10),
		UsdActiveFeedsCount: make(chan int, 10),
		EurActiveFeedsCount: make(chan int, 10),
	}

	ticker.SetFeeds([]BitcoinFeed{
		&Gdax{Currency: btcusd, Output: ticker.BtcUsd, FeedsCount: ticker.UsdActiveFeedsCount, Side: "buy"},
		&Btcc{Currency: btcusd, Output: ticker.BtcUsd, FeedsCount: ticker.UsdActiveFeedsCount, Side: "buy"},

		&Gdax{Currency: btceur, Output: ticker.BtcEur, FeedsCount: ticker.EurActiveFeedsCount, Side: "buy"},
	})

	for _, feed := range ticker.Feeds {
		go feed.Subscribe()
	}

	ticker.Print()
}
