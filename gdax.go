package main

import (
	"encoding/json"
	"strconv"
)

type BitcoinFeed interface {
	Subscribe()
	GetCurrency() string
}

var gdaxUrl = "wss://ws-feed.gdax.com"

type Gdax struct {
	Currency   string
	Output     chan float64
	FeedsCount chan int
	Side       string
}

type subscribeMessage struct {
	Type       string   `json:"type"`
	ProductIds []string `json:"product_ids"`
}

type response struct {
	Type      string `json:"type"`
	Price     string `json:"price"`
	ProductId string `json:"product_id"`
	Side      string `json:"side"`
}

func (feed *Gdax) Subscribe() {
	c := Conn{
		OnConnected: func(w *Conn) {
			feed.FeedsCount <- 1
			msg := subscribeMessage{"subscribe", []string{feed.Currency}}
			msgBytes, _ := json.Marshal(msg)
			w.Send(msgBytes)
		},
		OnMessage: func(msg []byte, w *Conn) {
			res := &response{}
			if err := json.Unmarshal(msg, res); err != nil {
				return
			}
			if res.Side != feed.Side {
				return
			}
			f, err := strconv.ParseFloat(res.Price, 64)
			if err == nil {
				feed.Output <- f
			}
		},
		OnError: func(err error) {
			feed.FeedsCount <- -1
		},
		Reconnect: true,
	}

	c.Dial(gdaxUrl, "")
}

func (feed *Gdax) GetCurrency() string {
	return feed.Currency
}
