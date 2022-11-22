package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// CheckBalance Basic usage of bitlaps API
// Using bitlaps https://developer.bitaps.com/blockchain which is free for 15 reqs in 5s currently
func CheckBalance(address string, fiat string) (float64, float64, string, error) {
	var (
		realBalance float64
		realPrice   float64
		extra       string
	)

	currency := "btc"
	multiplier := 0.00000001

	btcMatch, err := regexp.MatchString("^(1|3|bc1).", address)
	if err != nil {
		return 0, 0, "", err
	}

	if !btcMatch {
		ltcMatch, err := regexp.MatchString("(L|ltc1).", address)
		if err != nil {
			return 0, 0, "", err
		}

		if ltcMatch {
			currency = "ltc"
		} else if strings.HasPrefix(address, "0") {
			currency = "eth"
			multiplier = 0.000000000000000001
		} else if strings.HasPrefix(address, "D") {
			currency = "doge"
		} else {
			return 0, 0, "", errors.New("Unrecognised address. Must be Bitcoin, Litecoin or Ethereum")
		}
	}

	switch currency {
	case "btc":
		realBalance, err = mempoolGetBitcoinBalance(address, multiplier)
		if err != nil {
			realBalance, err = bitlapsGetBalance(currency, address, multiplier)
			if err != nil {
				return 0, 0, "", err
			}
		}

		realPrice, _ = coindeskGetBitcoinPrice(fiat, realBalance)
	case "eth":
		realBalance, realPrice, extra, err = ethplorerGetBalanceAndPrice(address)
		if err != nil {
			return 0, 0, "", err
		}

	case "ltc":
		realBalance, err = bitlapsGetBalance(currency, address, multiplier)
		if err != nil {
			return 0, 0, "", err
		}

		realPrice, _ = bitlapsGetPrice(currency, fiat, realBalance)

	case "doge":
		realBalance, realPrice, err = dogechainGetBalanceAndPrice(fiat, address)
		if err != nil {
			return 0, 0, "", err
		}
	}

	return realBalance, realPrice, extra, nil
}

func mempoolGetBitcoinBalance(address string, multiplier float64) (float64, error) {
	type balanceResp struct {
		Address    string `json:"address"`
		ChainStats struct {
			FundedTxoCount int   `json:"funded_txo_count"`
			FundedTxoSum   int64 `json:"funded_txo_sum"`
			SpentTxoCount  int   `json:"spent_txo_count"`
			SpentTxoSum    int64 `json:"spent_txo_sum"`
			TxCount        int   `json:"tx_count"`
		} `json:"chain_stats"`
		MempoolStats struct {
			FundedTxoCount int `json:"funded_txo_count"`
			FundedTxoSum   int `json:"funded_txo_sum"`
			SpentTxoCount  int `json:"spent_txo_count"`
			SpentTxoSum    int `json:"spent_txo_sum"`
			TxCount        int `json:"tx_count"`
		} `json:"mempool_stats"`
	}

	resultBytes, err := httpGet(fmt.Sprintf("https://mempool.space/api/address/%s", address))
	if err != nil {
		return 0, err
	}

	balance := balanceResp{}

	err = json.Unmarshal(resultBytes, &balance)
	if err != nil {
		return 0, err
	}

	realBalance := float64(balance.ChainStats.FundedTxoSum) * multiplier
	realBalance += float64(balance.MempoolStats.FundedTxoSum) * multiplier
	realBalance -= float64(balance.ChainStats.SpentTxoSum) * multiplier
	realBalance -= float64(balance.MempoolStats.SpentTxoSum) * multiplier

	return realBalance, nil
}

// coindeskGetBitcoinPrice get
func coindeskGetBitcoinPrice(fiat string, balance float64) (float64, error) {
	type bitcoinPriceResp struct {
		Time struct {
			Updated    string    `json:"updated"`
			UpdatedISO time.Time `json:"updatedISO"`
			Updateduk  string    `json:"updateduk"`
		} `json:"time"`
		Disclaimer string `json:"disclaimer"`
		ChartName  string `json:"chartName"`
		Bpi        struct {
			Usd struct {
				Code        string  `json:"code"`
				Symbol      string  `json:"symbol"`
				Rate        string  `json:"rate"`
				Description string  `json:"description"`
				RateFloat   float64 `json:"rate_float"`
			} `json:"USD"`
			Gbp struct {
				Code        string  `json:"code"`
				Symbol      string  `json:"symbol"`
				Rate        string  `json:"rate"`
				Description string  `json:"description"`
				RateFloat   float64 `json:"rate_float"`
			} `json:"GBP"`
			Eur struct {
				Code        string  `json:"code"`
				Symbol      string  `json:"symbol"`
				Rate        string  `json:"rate"`
				Description string  `json:"description"`
				RateFloat   float64 `json:"rate_float"`
			} `json:"EUR"`
		} `json:"bpi"`
	}

	resultBytes, err := httpGet("https://api.coindesk.com/v1/bpi/currentprice.json")
	if err != nil {
		return 0, err
	}

	bitcoinPrice := bitcoinPriceResp{}

	err = json.Unmarshal(resultBytes, &bitcoinPrice)
	if err != nil {
		return 0, err
	}

	fiat = strings.ToUpper(fiat)
	switch fiat {
	case "USD":
		return balance * bitcoinPrice.Bpi.Usd.RateFloat, nil
	case "GBP":
		return balance * bitcoinPrice.Bpi.Gbp.RateFloat, nil
	case "EUR":
		return balance * bitcoinPrice.Bpi.Eur.RateFloat, nil
	}

	return 0, fmt.Errorf("Price not found for fiat: %s", fiat)
}

// ethplorerGetBalanceAndPrice return value, price (usd only), extra string eg tokens including their value & price or return error
func ethplorerGetBalanceAndPrice(address string) (float64, float64, string, error) {
	const (
		currencySymbol = "$"                  // note: ethplorer returns USD only
		multiplier     = 0.000000000000000001 // TODO: any decimals != 18 will break
	)

	type ethplorerResp struct {
		Address string `json:"address"`
		Eth     struct {
			Price struct {
				Rate            float64 `json:"rate"`
				Diff            float64 `json:"diff"`
				Diff7D          float64 `json:"diff7d"`
				Ts              int     `json:"ts"`
				MarketCapUsd    float64 `json:"marketCapUsd"`
				AvailableSupply float64 `json:"availableSupply"`
				Volume24H       float64 `json:"volume24h"`
				VolDiff1        float64 `json:"volDiff1"`
				VolDiff7        float64 `json:"volDiff7"`
				VolDiff30       float64 `json:"volDiff30"`
				Diff30D         float64 `json:"diff30d"`
			} `json:"price"`
			Balance    float64 `json:"balance"`
			RawBalance string  `json:"rawBalance"`
		} `json:"ETH"`
		CountTxs int `json:"countTxs"`
		Tokens   []struct {
			TokenInfo struct {
				Address           string `json:"address"`
				Name              string `json:"name"`
				Decimals          string `json:"decimals"`
				Symbol            string `json:"symbol"`
				TotalSupply       string `json:"totalSupply"`
				Owner             string `json:"owner"`
				LastUpdated       int    `json:"lastUpdated"`
				IssuancesCount    int    `json:"issuancesCount"`
				HoldersCount      int    `json:"holdersCount"`
				Image             string `json:"image"`
				Description       string `json:"description"`
				Website           string `json:"website"`
				EthTransfersCount int    `json:"ethTransfersCount"`
				Price             struct {
					Rate            float64 `json:"rate"`
					Diff            float64 `json:"diff"`
					Diff7D          float64 `json:"diff7d"`
					Ts              int     `json:"ts"`
					MarketCapUsd    float64 `json:"marketCapUsd"`
					AvailableSupply int64   `json:"availableSupply"`
					Volume24H       float64 `json:"volume24h"`
					VolDiff1        float64 `json:"volDiff1"`
					VolDiff7        float64 `json:"volDiff7"`
					VolDiff30       float64 `json:"volDiff30"`
					Diff30D         float64 `json:"diff30d"`
					Bid             float64 `json:"bid"`
					Currency        string  `json:"currency"`
				} `json:"price"`
			} `json:"tokenInfo"`
			Balance    float64 `json:"balance"`
			TotalIn    int     `json:"totalIn"`
			TotalOut   int     `json:"totalOut"`
			RawBalance string  `json:"rawBalance"`
		} `json:"tokens"`
	}

	resultBytes, err := httpGet(fmt.Sprintf("https://api.ethplorer.io/getAddressInfo/%s?apiKey=freekey", address))
	if err != nil {
		return 0, 0, "", err
	}

	result := ethplorerResp{}

	err = json.Unmarshal(resultBytes, &result)
	if err != nil {
		return 0, 0, "", err
	}

	realBalance := result.Eth.Balance
	realPrice := result.Eth.Price.Rate * realBalance

	tokenListing := ""
	for idx := range result.Tokens {
		name := result.Tokens[idx].TokenInfo.Name
		symbol := result.Tokens[idx].TokenInfo.Symbol
		balance := float64(result.Tokens[idx].Balance) * multiplier
		price := balance * result.Tokens[idx].TokenInfo.Price.Rate

		tokenListing += fmt.Sprintf("\n- - %s (%s) %.08f %s%.02f", name, symbol, balance, currencySymbol, price)
	}

	return realBalance, realPrice, tokenListing, nil
}

func bitlapsGetBalance(currency, address string, multiplier float64) (float64, error) {
	// blockchainState holds the balance and some transaction infos
	type blockchainState struct {
		Data struct {
			Balance                  int    `json:"balance"`
			ReceivedAmount           int    `json:"receivedAmount"`
			ReceivedTxCount          int    `json:"receivedTxCount"`
			SentAmount               int    `json:"sentAmount"`
			SentTxCount              int    `json:"sentTxCount"`
			FirstReceivedTxPointer   string `json:"firstReceivedTxPointer"`
			FirstSentTxPointer       string `json:"firstSentTxPointer"`
			LastTxPointer            string `json:"lastTxPointer"`
			LargestReceivedTxAmount  int    `json:"largestReceivedTxAmount"`
			LargestReceivedTxPointer string `json:"largestReceivedTxPointer"`
			LargestSpentTxAmount     int    `json:"largestSpentTxAmount"`
			LargestSpentTxPointer    string `json:"largestSpentTxPointer"`
			ReceivedOutsCount        int    `json:"receivedOutsCount"`
			SpentOutsCount           int    `json:"spentOutsCount"`
			PendingReceivedAmount    int    `json:"pendingReceivedAmount"`
			PendingSentAmount        int    `json:"pendingSentAmount"`
			PendingReceivedTxCount   int    `json:"pendingReceivedTxCount"`
			PendingSentTxCount       int    `json:"pendingSentTxCount"`
			PendingReceivedOutsCount int    `json:"pendingReceivedOutsCount"`
			PendingSpentOutsCount    int    `json:"pendingSpentOutsCount"`
			Type                     string `json:"type"`
		} `json:"data"`
		Time float64 `json:"time"`
	}

	balanceBody, err := httpGet(fmt.Sprintf("https://api.bitaps.com/%s/v1/blockchain/address/state/%s", currency, address))
	if err != nil {
		return 0, err
	}

	balance := blockchainState{}
	err = json.Unmarshal(balanceBody, &balance)
	if err != nil {
		return 0, err
	}

	realBalance := float64(balance.Data.Balance) * multiplier

	return realBalance, nil
}

func bitlapsGetPrice(currency, fiat string, balance float64) (float64, error) {
	// ticker holds the currency ticker eg (btcusd)
	type ticker struct {
		Data struct {
			Last       float64 `json:"last"`
			LastChange float64 `json:"last_change"`
			Volume     float64 `json:"volume"`
			Open       float64 `json:"open"`
			High       float64 `json:"high"`
			Low        float64 `json:"low"`
			Markets    int     `json:"markets"`
		} `json:"data"`
		Time float64 `json:"time"`
	}

	priceBody, err := httpGet(fmt.Sprintf("https://api.bitaps.com/market/v1/ticker/%s%s", currency, fiat))
	if err != nil {
		return 0, err
	}

	price := ticker{}
	err = json.Unmarshal(priceBody, &price)
	if err != nil {
		return 0, err
	}

	realPrice := float64(balance) * price.Data.Last

	return realPrice, nil
}

func dogechainGetBalanceAndPrice(fiat, address string) (float64, float64, error) {
	type balanceData struct {
		Balance string `json:"balance"`
		Success int
	}

	type priceData struct {
		Status string `json:"status"`
		Data   struct {
			Network string `json:"network"`
			Prices  []struct {
				Price     string `json:"price"`
				PriceBase string `json:"price_base"`
				Exchange  string `json:"exchange"`
				Time      int    `json:"time"`
			} `json:"prices"`
		} `json:"data"`
	}

	balanceBytes, err := httpGet("https://dogechain.info/api/v1/address/balance/" + address)
	if err != nil {
		return 0, 0, err
	}

	balanceBody := balanceData{}
	err = json.Unmarshal(balanceBytes, &balanceBody)
	if err != nil {
		return 0, 0, err
	}

	balanceFloat, err := strconv.ParseFloat(balanceBody.Balance, 64)
	if err != nil {
		return 0, 0, err
	}

	priceBytes, err := httpGet("https://sochain.com/api/v2/get_price/DOGE/USD")
	if err != nil {
		return 0, 0, err
	}

	priceBody := priceData{}
	err = json.Unmarshal(priceBytes, &priceBody)
	if err != nil {
		return 0, 0, err
	}

	priceFloat, err := strconv.ParseFloat(priceBody.Data.Prices[0].Price, 64)
	if err != nil {
		return 0, 0, err
	}

	return balanceFloat, priceFloat * balanceFloat, nil
}

func httpGet(url string) ([]byte, error) {
	client := http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	return io.ReadAll(res.Body)
}
