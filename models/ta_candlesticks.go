package models

type TACadlestick struct {
	OpenPrice  float64
	HighPrice  float64
	ClosePrice float64
	Ema223     float64
	Ema20      float64
	RSI14      float64
}

// This will tell if in the single point we are in bull or in bear
// Equals situation are considered Bear so false will return
func (tac *TACadlestick) IsBull() bool {
	return tac.Ema20 > tac.Ema223
}

// In this case as normal interval we consider A <= RSI14 <= B this is considered
//
// With A(bottomLimit) > B(topLimit).
//
// As "normal" interval for a strategy.
//
// But since A>B and B can be considered as max value, we can say that if we go over
//
// A then the price started to be in a "normal" interval
//
// For this reason we will have also another params:
//
// @params bottomLimit Lower limit for the RSI
//
// @params topLimit Upper limit for the RSI
//
// @params strict It asses if we have to strictly check for the RSI in that interval
func (tac *TACadlestick) RSIInNormalInterval(bottomLimit float64, topLimit float64, strict bool) bool {
	if !strict {
		return tac.RSI14 >= bottomLimit && tac.RSI14 <= bottomLimit
	}
	return tac.RSI14 >= bottomLimit
}

func (tac *TACadlestick) IsGreen() bool {
	return tac.ClosePrice > tac.OpenPrice
}
