package trader

import (
	"core_trade_bot/models"

	bc "github.com/monkeymatt0/binance_client"
)

// The trader based on the market trand can use a specific strategy to buy and sell
// Trigger function is used to understand if some trigger condition has happened
type Trader interface {
	OpportunityFound() bool
	Trigger(candelsticks []bc.RawCandlestick)
	BuyTechnique(strategy models.BuyStrategy, trend models.MarketTrend) (uint64, error)
	SellTechnique(take_profit float64, stop_loss float64)
	CancelOrderParams(params map[string]string) // Function will cancel the order if the order does not fire on time
	CheckOrder(listenKey string)                // This function is used to check the order status
}
