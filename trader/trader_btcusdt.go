package trader

import (
	"core_trade_bot/models"
	"core_trade_bot/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	talib "github.com/markcheno/go-talib"
	bc "github.com/monkeymatt0/binance_client"
)

type BtcUsdtTrader struct {
	apiKey    string
	secretKey string

	usdt float64

	buyPrice    float64 // Price I bought
	tpSellPrice float64 // Price I sell in profit
	slSellPrice float64 // Price I loss in profit

	triggerFired    bool
	buyOrderPlaced  bool
	sellOrderPlaced bool
	stopLoss        float64 // optimal value 0.004
	stopPriceLoss   float64 // optimal value 0.003
	takeProfit      float64 // optimal value 0.025

	bottomLimitRSI float64
	topLimitRSI    float64

	validCandle uint64        //optimal value 5 This value will tell how many candle should the order last, if in validity candle the order does not took place then the order is invalidate
	stopHours   time.Duration //optimal value 2 This will tell how many hours wi should stop
	binance     bc.Binance
}

// This New method will set the proper value for the attribute:
// @param test bool => if true we are using a test network
// @param apiKey string => It's the API of exchange
// @param secretKey string => It's the secretKey of the api
// @param bl float64 => Is the bottom limit for the RSI
// @param tl float64 => Is the top limit for the RSI
// @param tp float64 => Is the take profit of the strategy for a sell operation
// @param sl float64 => Is the stop loss used for a sell operation
// @param validCandle uint8 => It says for how many candle the order is valid
// @param stopHours time.Duration => It says asses the amount of time the bot has to be stopped before it restart operate
func (but *BtcUsdtTrader) New(
	test bool,
	apiKey, secretKey string,
	bl, tl, tp, sl, spl float64,
	validCandle uint64,
	stopHours time.Duration,
) error {
	but.binance.New(test)
	but.apiKey = apiKey
	but.secretKey = secretKey
	but.bottomLimitRSI = bl
	but.topLimitRSI = tl
	but.stopHours = stopHours
	but.validCandle = validCandle
	but.stopLoss = sl
	but.stopPriceLoss = spl
	but.takeProfit = tp
	usdt, err := but.usdtInWallet(models.USDT)
	if err != nil {
		return err
	}
	but.usdt = usdt
	return nil
}

func (but *BtcUsdtTrader) getFetchParams() map[string]string {
	ret := make(map[string]string)
	ret["interval"] = "5m"
	ret["symbol"] = strings.Join([]string{string(models.BTC), string(models.USDT)}, "")
	ret["limit"] = "250"
	return ret
}

func (but *BtcUsdtTrader) OpportunityFound() (bool, error) {
	params := but.getFetchParams()
	candlesticks, err := but.binance.KlinesRequest(params)
	if err != nil {
		return false, nil
	}
	cCandlesticks := utils.CloseCandlesticks(candlesticks)
	oCandlesticks := utils.OpenCandlesticks(candlesticks)
	hCandlesticks := utils.HighCandlesticks(candlesticks)
	ema223 := talib.Ema(cCandlesticks, 223)
	ema20 := talib.Ema(cCandlesticks, 20)
	rsi := talib.Rsi(cCandlesticks, 14)

	TACandlesticks := utils.CreateTACandlesticks(oCandlesticks, cCandlesticks, hCandlesticks, ema223, ema20, rsi)
	lastThreeCandles := TACandlesticks[len(TACandlesticks)-3:]
	but.Trigger(lastThreeCandles)
}

// This method will fetch the current amount of USDT that are in the wallet
func (but *BtcUsdtTrader) usdtInWallet(coin models.Coin) (float64, error) {
	accountInfo, err := but.binance.AccountRequest(make(map[string]string), but.apiKey, but.secretKey)
	if err != nil {
		return 0.0, err
	}
	for _, value := range accountInfo.Balances {
		if value.Asset == string(coin) {
			return strconv.ParseFloat(value.Free, 64)
		}
	}
	return 0.0, nil // @todo : Here I should add a custom error
}

// Trigger checks if a specific trigger happened:
// @param lastThreeCandles []models.TACandlestick => Are the last 3 candles from the candlstick chart with tachnicals indicator
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
	orderIds, err := but.binance.OrderRequest(params, but.secretKey, but.secretKey, http.MethodPost)
	if err != nil {
		return 0, err
	}
	but.buyOrderPlaced = true
	listenKey, err := but.getListenKey()
	if err != nil {
		return 0, err
	}
	but.CheckOrder(listenKey)
	return orderIds[0], nil
	// @todo : Once this is done you should periodically check that the order
	// is executed within the valid_candle period.
}

func (but *BtcUsdtTrader) buyParamsSetter(lastCandlestick models.TACadlestick) map[string]string {
	ret := make(map[string]string)
	but.buyPrice = (lastCandlestick.Ema223 - lastCandlestick.Ema20) / 2
	ret["symbol"] = "BTCUSDT"
	ret["side"] = "BUY"
	ret["type"] = "STOP_LOSS_LIMIT"
	ret["timeInForce"] = "GTC"
	ret["stopPrice"] = strconv.FormatFloat((but.buyPrice)-1, 'f', 6, 64)
	ret["price"] = strconv.FormatFloat((but.buyPrice)/2, 'f', 6, 64)
	ret["quantity"] = strconv.FormatFloat(but.usdt/but.buyPrice, 'f', 6, 64)
	ret["timestamp"] = strconv.FormatInt(time.Now().Unix()*1000, 10)
	return ret
}

func (but *BtcUsdtTrader) sellParamsSetter() map[string]string {
	ret := make(map[string]string)
	qty, err := but.usdtInWallet(models.BTC)
	if err != nil {
		return ret
	}

	abovePrice := but.buyPrice + (but.buyPrice * but.takeProfit)
	belowStopPrice := but.buyPrice - (but.buyPrice * but.stopPriceLoss)
	belowPrice := but.buyPrice - (but.buyPrice * but.stopLoss)
	ret["symbol"] = "BTCUSDT"
	ret["side"] = "SELL"
	ret["quantity"] = strconv.FormatFloat(qty, 'f', 6, 64)
	ret["aboveType"] = "LIMIT_MAKER"
	ret["abovePrice"] = strconv.FormatFloat(abovePrice, 'f', 6, 64)
	ret["belowType"] = "STOP_LOSS_LIMIT"
	ret["belowPrice"] = strconv.FormatFloat(belowPrice, 'f', 6, 64)
	ret["belowStopPrice"] = strconv.FormatFloat(belowStopPrice, 'f', 6, 64)
	ret["belowTimeInForce"] = "GTC"
	ret["timestamp"] = strconv.FormatInt(time.Now().Unix()*1000, 10)
	return ret
}

func (but *BtcUsdtTrader) getListenKey() (*string, error) {
	key, err := but.binance.ListenKeyRequest(but.apiKey, 0)
	if err != nil {
		return nil, err
	}
	return &key.ListenKey, err
}

func (but *BtcUsdtTrader) SellTechnique() (*[]uint64, error) {
	sellParams := but.sellParamsSetter()
	orderIds, err := but.binance.OrderRequest(sellParams, but.secretKey, but.secretKey, http.MethodPost)
	if err != nil {
		return nil, err
	}
	but.sellOrderPlaced = true
	return &orderIds, nil
}

func (but *BtcUsdtTrader) CheckOrder(listenKey *string) {
	// @todo : Use the web socket stream to have update regarding the order
	if listenKey == nil {
		return
	}
	validTime := but.validCandle * 5 * uint64(time.Minute)
	but.binance.UserDataStreamSocket(*listenKey, time.Duration(validTime))
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

	// Check to see if the candle is below the RSI bottom
	if candlesticks[0].RSI14 > bottomLimit {
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
