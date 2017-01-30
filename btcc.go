package main

import (
	"encoding/json"
	"strconv"
	"strings"
)

var btccUrl = "wss://api.bitfinex.com/ws"

type Btcc struct {
	Currency   string
	Output     chan float64
	FeedsCount chan int
	Side       string
}

type btccSubscribeMessage struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
	Pair    string `json:"pair"`
}

func (feed *Btcc) Subscribe() {
	c := Conn{
		OnConnected: func(w *Conn) {
			feed.FeedsCount <- 1
			msg := btccSubscribeMessage{"subscribe", "ticker", feed.LocalCurrency()}
			msgBytes, _ := json.Marshal(msg)
			w.Send(msgBytes)
		},
		OnMessage: func(msg []byte, w *Conn) {
			items := strings.Split(string(msg), ",")
			if len(items) < 10 {
				return
			}
			f, err := strconv.ParseFloat(items[1], 64)
			if err == nil {
				feed.Output <- f
			}
		},
		OnError: func(err error) {
			feed.FeedsCount <- -1
		},
		Reconnect: true,
	}

	c.Dial(btccUrl, "")
}

func (feed *Btcc) GetCurrency() string {
	return feed.Currency
}

func (feed *Btcc) LocalCurrency() string {
	switch {
	case feed.Currency == btcusd:
		return "BTCUSD"
	case feed.Currency == btceur:
		return "BTCEUR"
	}
	return feed.Currency
}
