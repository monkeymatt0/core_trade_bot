package models

type Config struct {
	ApiKey           string  `json:"api_key"`
	SecretKey        string  `json:"secret_key"`
	RsiBottomLimit   float64 `json:"rsi_bottom_limit"`
	RsiTopLimit      float64 `json:"rsi_top_limit"`
	TakeProfit       float64 `json:"take_profit"`
	StopLoss         float64 `json:"stop_loss"`
	StopPriceLoss    float64 `json:"stop_price_loss"`
	StopHours        uint8   `json:"stop_hours"`
	ValidCandle      uint64  `json:"valid_candle"`
	MaxCapitalUsable float64 `json:"max_capital_usable"`
}
