package main

import (
	"encoding/json"
	"strconv"
)

var gdaxUrl = "wss://ws-feed.gdax.com"

type Gdax struct {
	Currency   string
	Output     chan float64
	FeedsCount chan int
	Side       string
}

type gdaxSubscribeMessage struct {
	Type       string   `json:"type"`
	ProductIds []string `json:"product_ids"`
}

type gdaxResponse struct {
	Type      string `json:"type"`
	Price     string `json:"price"`
	ProductId string `json:"product_id"`
	Side      string `json:"side"`
}

func (feed *Gdax) Subscribe() {
	c := Conn{
		OnConnected: func(w *Conn) {
			feed.FeedsCount <- 1
			msg := gdaxSubscribeMessage{"subscribe", []string{feed.GetCurrency()}}
			msgBytes, _ := json.Marshal(msg)
			w.Send(msgBytes)
		},
		OnMessage: func(msg []byte, w *Conn) {
			res := &gdaxResponse{}
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
