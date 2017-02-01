package main

import "testing"

func TestTicker_SetFeeds(t *testing.T) {
	ticker := Ticker{
		BtcUsd:              make(chan float64, 10),
		BtcEur:              make(chan float64, 10),
		UsdActiveFeedsCount: make(chan int, 10),
		EurActiveFeedsCount: make(chan int, 10),
	}

	ticker.SetFeeds([]BitcoinFeed{
		&Gdax{Currency: btcusd, Output: ticker.BtcUsd, FeedsCount: ticker.UsdActiveFeedsCount, Side: "buy"},
		&Gdax{Currency: btcusd, Output: ticker.BtcUsd, FeedsCount: ticker.UsdActiveFeedsCount, Side: "buy"},

		&Gdax{Currency: btceur, Output: ticker.BtcEur, FeedsCount: ticker.EurActiveFeedsCount, Side: "buy"},
	})

	if ticker.UsdFeedsCount != 2 {
		t.Errorf("UsdFeedsCount expected = 2, got = '%s'", ticker.UsdFeedsCount)
	}
	if ticker.EurFeedsCount != 1 {
		t.Errorf("EurFeedsCount expected = 1, got = '%s'", ticker.EurFeedsCount)
	}
}
