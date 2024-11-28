package main

import (
	"core_trade_bot/models"
	"core_trade_bot/trader"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func main() {
	config := models.Config{}
	bytes, err := os.ReadFile("testConfig.json")
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
		config.StopPriceLoss,
		config.ValidCandle, // Candles validity for order
		stopHoursDuration,  // Hours of stop after a profit/loss stop
	)

	huntingCh := make(chan struct{})        // When this channel will receive a signal then a check of the candlestick will take place
	placeBuyOrderCh := make(chan struct{})  // When this channel will receive an order a buy order will be placed
	checkOrderCh := make(chan string)       // This channel will receive the listen key
	placeSellOrderCh := make(chan struct{}) // When this channel receive a signal, a sell order will be placed

	// @remind : macro area to develop
	// 1) Hunting (channel) phase - In this phase the trader is looking for opportunity
	// 2) Place Order (Buy - channel) - During this phase the trader if an opportunity arise place a buy order
	// 3) Check (channel) the order - During this phase the trader check the order execution and if too much time passes, he will close the order
	// 4) Place Order (Sell - channel) - During this phase the trader will place the sell order
	// 5) Check (channel) the order - During this phase the trader will check the order, when it will be executed, start again from 1) Hunting
}

// package main

// import (
// 	"crypto/hmac"
// 	"crypto/sha256"
// 	"encoding/hex"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"net/url"
// 	"strconv"
// 	"time"
// )

// const (
// 	apiKey    = "0PM8UeDy2GRvQ94XVIJhnZAC8zmrJg3kD32TQGyWjwFrxMktJdM0VFfH2jJDIMIl" // Replace with your Binance API Key
// 	apiSecret = "7MvD1jmpLjRiYqftER72SVehI4J4zppIOZqRcdwjCv0ha1Y4IAinGazG4rYH2ERg" // Replace with your Binance API Secret
// 	baseURL   = "https://testnet.binance.vision"
// )

// // Sign the query string using HMAC SHA256
// func sign(query string) string {
// 	mac := hmac.New(sha256.New, []byte(apiSecret))
// 	mac.Write([]byte(query))
// 	return hex.EncodeToString(mac.Sum(nil))
// }

// // Place a Market Order
// func placeMarketOrder(symbol, side, quantity string) {
// 	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
// 	params := url.Values{}
// 	params.Add("symbol", symbol)     // Trading pair, e.g., "BTCUSDT"
// 	params.Add("side", side)         // "BUY" or "SELL"
// 	params.Add("type", "MARKET")     // Market order type
// 	params.Add("quantity", quantity) // Amount to buy or sell
// 	params.Add("timestamp", timestamp)

// 	signature := sign(params.Encode())
// 	params.Add("signature", signature)

// 	reqURL := fmt.Sprintf("%s/api/v3/order", baseURL)
// 	req, err := http.NewRequest("POST", reqURL, nil)
// 	if err != nil {
// 		panic(err)
// 	}
// 	req.Header.Set("X-MBX-APIKEY", apiKey)
// 	req.URL.RawQuery = params.Encode()

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer resp.Body.Close()

// 	body, _ := io.ReadAll(resp.Body)
// 	fmt.Println("Response from placing market order:", string(body))
// }

// // Place an OCO Order
// func placeOCOOrder() int64 {
// 	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
// 	params := url.Values{}
// 	params.Add("symbol", "BTCUSDT")
// 	params.Add("side", "SELL")
// 	params.Add("quantity", "0.001") // Adjust quantity
// 	params.Add("aboveType", "LIMIT_MAKER")
// 	params.Add("belowType", "STOP_LOSS_LIMIT")
// 	params.Add("abovePrice", "96000") // Take-profit price
// 	params.Add("belowPrice", "90000")
// 	params.Add("belowStopPrice", "89999")
// 	params.Add("belowTimeInForce", "GTC")
// 	params.Add("timestamp", timestamp)

// 	signature := sign(params.Encode())
// 	params.Add("signature", signature)

// 	reqURL := fmt.Sprintf("%s/api/v3/orderList/oco", baseURL)
// 	req, err := http.NewRequest("POST", reqURL, nil)
// 	if err != nil {
// 		panic(err)
// 	}
// 	req.Header.Set("X-MBX-APIKEY", apiKey)
// 	req.URL.RawQuery = params.Encode()

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer resp.Body.Close()

// 	body, _ := io.ReadAll(resp.Body)
// 	fmt.Println("Response from placing OCO Order:", string(body))

// 	// Simulated response parsing (adapt for actual API)
// 	var orderID int64 = 12345678 // Replace with parsed order ID
// 	return orderID
// }

// // Check Order Status
// func checkOrderStatus(orderID int64) {
// 	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
// 	params := url.Values{}
// 	params.Add("symbol", "BTCUSDT")
// 	params.Add("orderId", strconv.FormatInt(orderID, 10))
// 	params.Add("timestamp", timestamp)

// 	signature := sign(params.Encode())
// 	params.Add("signature", signature)

// 	reqURL := fmt.Sprintf("%s/api/v3/order", baseURL)
// 	req, err := http.NewRequest("GET", reqURL, nil)
// 	if err != nil {
// 		panic(err)
// 	}
// 	req.Header.Set("X-MBX-APIKEY", apiKey)
// 	req.URL.RawQuery = params.Encode()

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer resp.Body.Close()

// 	body, _ := io.ReadAll(resp.Body)
// 	fmt.Println("Response from checking order status:", string(body))
// }

// func main() {
// 	// // Example usage: Place a market buy order for 0.01 BTC in BTCUSDT
// 	// fmt.Println("Placing Market Order...")
// 	// placeMarketOrder("BTCUSDT", "BUY", "0.01")

// 	// // Step 1: Place an OCO Order
// 	// fmt.Println("Placing OCO Order...")
// 	// placeOCOOrder()
// 	orderID := 7910687

// 	// Step 2: Check Order Status
// 	fmt.Println("Checking Order Status...")
// 	checkOrderStatus(int64(orderID))
// }
