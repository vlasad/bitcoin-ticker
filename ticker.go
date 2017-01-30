package main

import "fmt"

type Ticker struct {
	BtcUsd              chan float64
	BtcEur              chan float64
	UsdActiveFeedsCount chan int
	EurActiveFeedsCount chan int
	Feeds               []BitcoinFeed
	UsdFeedsCount       int
	EurFeedsCount       int
}

func (t *Ticker) Print() {
	var _btcusd, _btceur float64
	var _usdSourcesCount, _eurSourcesCount int
	format := "BTC/USD: %6.2f   EUR/USD: %4.2f   BTC/EUR: %6.2f   Active sources: BTC/USD (%d of %d) EUR/USD (%d of %d)\n"

	for {
		select {
		case x := <-t.UsdActiveFeedsCount:
			_usdSourcesCount += x
		case x := <-t.EurActiveFeedsCount:
			_eurSourcesCount += x
		case _btcusd = <-t.BtcUsd:
			print("\033[H\033[2J")
			fmt.Printf(format, _btcusd, _btcusd/_btceur, _btceur, _usdSourcesCount, t.UsdFeedsCount, _eurSourcesCount, t.EurFeedsCount)
		case _btceur = <-t.BtcEur:
			print("\033[H\033[2J")
			fmt.Printf(format, _btcusd, _btcusd/_btceur, _btceur, _usdSourcesCount, t.UsdFeedsCount, _eurSourcesCount, t.EurFeedsCount)
		}
	}
}

func (t *Ticker) CountUsdSources() {
	t.UsdFeedsCount = 0
	for _, feed := range t.Feeds {
		if feed.GetCurrency() == btcusd {
			t.UsdFeedsCount += 1
		}
	}
}

func (t *Ticker) CountEurSources() {
	t.EurFeedsCount = 0
	for _, feed := range t.Feeds {
		if feed.GetCurrency() == btceur {
			t.EurFeedsCount += 1
		}
	}
}

func (t *Ticker) SetFeeds(feeds []BitcoinFeed) {
	t.Feeds = feeds
	t.CountUsdSources()
	t.CountEurSources()
}
