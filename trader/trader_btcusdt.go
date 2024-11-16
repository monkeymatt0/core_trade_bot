package trader

import (
	"core_trade_bot/models"
	"net/http"
	"strconv"
	"time"

	bc "github.com/monkeymatt0/binance_client"
)

type BtcUsdtTrader struct {
	apiKey    string
	secretKey string

	usdt float64

	triggerFired    bool
	buyOrderPlaced  bool
	sellOrderPlaced bool
	stopLoss        float64 // optimal value 0.003
	takeProfit      float64 // optimal value 0.025

	bottomLimitRSI float64
	topLimitRSI    float64

	validCandle uint8         //optimal value 5 This value will tell how many candle should the order last, if in validity candle the order does not took place then the order is invalidate
	stopHours   time.Duration //optimal value 2 This will tell how many hours wi should stop
	binance     bc.Binance
}

// @todo : Add all the proper setting for the trader
func (but *BtcUsdtTrader) New(test bool, apiKey, secretKey string, bl, tl, tp, sl float64, validCandle uint8, stopHours time.Duration) error {
	but.binance.New(test)
	but.apiKey = apiKey
	but.secretKey = secretKey
	but.bottomLimitRSI = bl
	but.topLimitRSI = tl
	but.stopHours = stopHours
	but.validCandle = validCandle
	but.stopLoss = sl
	but.takeProfit = tp
	usdt, err := but.usdtInWallet()
	if err != nil {
		return err
	}
	but.usdt = usdt
	return nil
}

func (but *BtcUsdtTrader) usdtInWallet() (float64, error) {
	accountInfo, err := but.binance.AccountRequest(make(map[string]string), but.apiKey, but.secretKey)
	if err != nil {
		return 0.0, err
	}
	for _, value := range accountInfo.Balances {
		if value.Asset == "USDT" {
			return strconv.ParseFloat(value.Free, 64)
		}
	}
	return 0.0, nil // @todo : Here I should add a custom error
}

func (but *BtcUsdtTrader) Trigger(lastThreeCandles []models.TACadlestick) {
	if !but.buyOrderPlaced &&
		!lastThreeCandles[0].IsBull() &&
		!lastThreeCandles[1].IsBull() &&
		!lastThreeCandles[2].IsBull() &&
		but.monkeyTrigger(lastThreeCandles, but.bottomLimitRSI, but.topLimitRSI) {

		but.triggerFired = true
		but.BuyTechnique(models.MKY_IVN, models.BEAR, lastThreeCandles[2])
	}
}

// BuyTechnique will perform the buy and will set
func (but *BtcUsdtTrader) BuyTechnique(strategy models.BuyStrategy, trend models.MarketTrend, lastCandlestick models.TACadlestick) (uint64, error) {
	params := but.buyParamsSetter(lastCandlestick)
	orderId, err := but.binance.OrderRequest(params, but.secretKey, but.secretKey, http.MethodPost)
	if err != nil {
		return 0, err
	}
	but.buyOrderPlaced = true
	return orderId, nil
	// @todo : Once this is done you should periodically check that the order
	// is executed within the valid_candle period.
}

func (but *BtcUsdtTrader) buyParamsSetter(lastCandlestick models.TACadlestick) map[string]string {
	ret := make(map[string]string)
	ret["symbol"] = "BTCUSDT"
	ret["side"] = "BUY"
	ret["type"] = "STOP_LOSS_LIMIT"
	ret["timeInForce"] = "GTC"
	ret["stopPrice"] = strconv.FormatFloat(((lastCandlestick.Ema223-lastCandlestick.Ema20)/2)-1, 'f', 6, 64)
	ret["price"] = strconv.FormatFloat((lastCandlestick.Ema223-lastCandlestick.Ema20)/2, 'f', 6, 64)
	ret["quantity"] = strconv.FormatFloat(but.usdt, 'f', 6, 64)
	ret["timestamp"] = strconv.FormatInt(time.Now().Unix()*1000, 10)
	return ret
}

func (but *BtcUsdtTrader) CheckOrder(listenKey string) {
	// @todo : Use the web socket stream to have update regarding the order
}

// @todo : implement SellTechnique
func (but *BtcUsdtTrader) SellTechnique(params map[string]string) {

	but.sellOrderPlaced = true
	// @todo : Once the sell order is placed you should check if it is executed
}

// moonketTrigger function check the following:
//
// RSI dropped under the 45 value in the last 3 candle
// 2 green candle follwed this drop
// these 2 green candle make RSI return in the normal interval

func (but *BtcUsdtTrader) monkeyTrigger(candlesticks []models.TACadlestick, bottomLimit float64, topLimit float64) bool {

	if len(candlesticks) > 3 || len(candlesticks) < 3 {
		return false
	}

	// Check to see if the candle is below the 45 RSI
	if candlesticks[0].RSI14 > 45 {
		return false
	}

	// Check to see if the next 2 candles are green
	if !candlesticks[1].IsGreen() && !candlesticks[2].IsGreen() {
		return false
	}

	// Check to see if the next 2 candles returned in the normal interval
	if !candlesticks[1].RSIInNormalInterval(bottomLimit, topLimit, false) && !candlesticks[2].RSIInNormalInterval(bottomLimit, topLimit, false) {
		return false
	}

	// If none of this false is returned then the condition happened and the trigger is fired
	return true
}
