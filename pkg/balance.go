package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

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

// CheckBalance Basic usage of bitlaps API
// Using bitlaps https://developer.bitaps.com/blockchain which is free for 15 reqs in 5s currently
func CheckBalance(address string, fiat string) (float64, float64, error) {
	currency := "btc"
	multiplier := 0.00000001

	btcMatch, err := regexp.MatchString("(1|3|bc1).", address)
	if err != nil {
		return 0, 0, err
	}

	if !btcMatch {
		ltcMatch, err := regexp.MatchString("(L|ltc1).", address)
		if err != nil {
			return 0, 0, err
		}
		if ltcMatch {
			currency = "ltc"
		} else if strings.HasPrefix(address, "0") {
			currency = "eth"
			multiplier = 0.000000000000000001
		} else {
			return 0, 0, errors.New("Unrecognised address. Must be Bitcoin, Litecoin or Ethereum")
		}
	}

	balanceBody, err := httpGet(fmt.Sprintf("https://api.bitaps.com/%s/v1/blockchain/address/state/%s", currency, address))
	if err != nil {
		return 0, 0, err
	}

	balance := blockchainState{}
	err = json.Unmarshal(balanceBody, &balance)
	if err != nil {
		return 0, 0, err
	}

	priceBody, err := httpGet(fmt.Sprintf("https://api.bitaps.com/market/v1/ticker/%s%s", currency, fiat))
	if err != nil {
		return 0, 0, err
	}

	price := ticker{}
	err = json.Unmarshal(priceBody, &price)
	if err != nil {
		return 0, 0, err
	}

	realBalance := float64(balance.Data.Balance) * multiplier
	realPrice := float64(realBalance) * price.Data.Last

	return realBalance, realPrice, nil
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

	return ioutil.ReadAll(res.Body)
}
