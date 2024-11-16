package utils

import (
	"core_trade_bot/models"

	bc "github.com/monkeymatt0/binance_client"
)

func CloseCandlesticks(candlesticks []bc.RawCandlestick) []float64 {
	closes := []float64{}

	for _, value := range candlesticks {
		closes = append(closes, value.ClosePrice)
	}
	return closes
}

func HighCandlesticks(candlesticks []bc.RawCandlestick) []float64 {
	Highs := []float64{}

	for _, value := range candlesticks {
		Highs = append(Highs, value.HighPrice)
	}
	return Highs
}

func OpenCandlesticks(candlesticks []bc.RawCandlestick) []float64 {
	Opens := []float64{}

	for _, value := range candlesticks {
		Opens = append(Opens, value.OpenPrice)
	}
	return Opens
}

func CreateTACandlesticks(opens []float64, closes []float64, highs []float64, ema223 []float64, ema20 []float64, rsi []float64) []models.TACadlestick {
	taCandlesticks := []models.TACadlestick{}

	for index := 0; index < len(closes); index++ {
		taCandlesticks = append(taCandlesticks, models.TACadlestick{
			OpenPrice:  opens[index],
			HighPrice:  highs[index],
			ClosePrice: closes[index],
			Ema223:     ema223[index],
			Ema20:      ema20[index],
			RSI14:      rsi[index],
		})
	}
	return taCandlesticks
}
