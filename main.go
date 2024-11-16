package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"core_trade_bot/models"
	"core_trade_bot/trader"
)

func main() {
	config := models.Config{}
	bytes, err := os.ReadFile("config.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(bytes))
	if err := json.Unmarshal(bytes, &config); err != nil {
		fmt.Println(err)
	}

	btcusdtTrader := &trader.BtcUsdtTrader{}
	stopHoursDuration := time.Duration(config.StopHours) * time.Hour
	btcusdtTrader.New(
		true,
		config.ApiKey,
		config.SecretKey,
		config.RsiBottomLimit, // Low limit for the RSI
		config.RsiTopLimit,    // Top Limit for the RSI
		config.TakeProfit,
		config.StopLoss,
		config.ValidCandle, // Candles validity for order
		stopHoursDuration,  // Hours of stop after a profit/loss stop
	)
}
