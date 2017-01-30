package main

type BitcoinFeed interface {
	Subscribe()
	GetCurrency() string
}
