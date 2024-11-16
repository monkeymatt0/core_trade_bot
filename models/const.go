package models

type BuyStrategy uint8

const (
	MKY BuyStrategy = iota
	MKY_IVN
)

type MarketTrend uint8

const (
	BEAR MarketTrend = iota
	BULL
)
