package internal

import (
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestCheckBalance(t *testing.T) {
	type args struct {
		address string
		fiat    string
	}

	const (
		mempoolBitcoinResp       = `{"address":"1FHxL2JskCy6g98wEMxpaNNkxohjq3hUKk","chain_stats":{"funded_txo_count":1,"funded_txo_sum":1010000,"spent_txo_count":0,"spent_txo_sum":0,"tx_count":1},"mempool_stats":{"funded_txo_count":0,"funded_txo_sum":0,"spent_txo_count":0,"spent_txo_sum":0,"tx_count":0}}`
		coindeskBitcoinPriceResp = `{"time":{"updated":"Jul 3, 2022 09:39:00 UTC","updatedISO":"2022-07-03T09:39:00+00:00","updateduk":"Jul 3, 2022 at 10:39 BST"},"disclaimer":"This data was produced from the CoinDesk Bitcoin Price Index (USD). Non-USD currency data converted using hourly conversion rate from openexchangerates.org","chartName":"Bitcoin","bpi":{"USD":{"code":"USD","symbol":"&#36;","rate":"18,974.5602","description":"United States Dollar","rate_float":18974.5602},"GBP":{"code":"GBP","symbol":"&pound;","rate":"15,674.9791","description":"British Pound Sterling","rate_float":15674.9791},"EUR":{"code":"EUR","symbol":"&euro;","rate":"18,193.8709","description":"Euro","rate_float":18193.8709}}}`
		ethplorerResp            = `{"address":"0x76270d9d9afc0cf4ebffbafe6401e01cb0f021ce","ETH":{"price":{"rate":1056.0434559974965,"diff":1.89,"diff7d":-15.05,"ts":1656841080,"marketCapUsd":128194943331.71054,"availableSupply":121391731.1865,"volume24h":9216735519.421448,"volDiff1":-36.1791554565028,"volDiff7":-8.91690879661894,"volDiff30":-11.340592431526801,"diff30d":-41.48069018787955},"balance":0.123,"rawBalance":"123000000000000000"},"countTxs":1,"tokens":[{"tokenInfo":{"address":"0xae78736cd615f374d3085123a210448e74fc6393","decimals":"18","name":"Rocket Pool ETH","symbol":"RETH","totalSupply":"95756543809751930313988","lastUpdated":1656837903,"issuancesCount":6186,"holdersCount":4123,"website":"https://rocketpool.net","image":"/images/RETHae78736c.png","ethTransfersCount":0,"price":{"rate":1077.5038815109424,"diff":2.42,"diff7d":-14.73,"ts":1656840960,"marketCapUsd":0,"availableSupply":0,"volume24h":451552.9133159675,"volDiff1":186.69939171602368,"volDiff7":5.438980995111038,"volDiff30":-99.97352203631957,"diff30d":-50.01024247936624,"bid":2946.28,"currency":"USD"}},"balance":1.234e+17,"totalIn":0,"totalOut":0,"rawBalance":"123400000000000000"}]}`
	)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://mempool.space/api/address/1FHxL2JskCy6g98wEMxpaNNkxohjq3hUKk",
		httpmock.NewStringResponder(200, mempoolBitcoinResp))
	httpmock.RegisterResponder("GET", "https://api.coindesk.com/v1/bpi/currentprice.json",
		httpmock.NewStringResponder(200, coindeskBitcoinPriceResp))
	httpmock.RegisterResponder("GET", "https://api.ethplorer.io/getAddressInfo/0x76270d9D9afC0cf4EbfFBafE6401E01cb0F021Ce?apiKey=freekey",
		httpmock.NewStringResponder(200, ethplorerResp))

	tests := []struct {
		name    string
		args    args
		want    float64
		want1   float64
		wantErr bool
	}{
		{name: "valid btc", args: args{address: "1FHxL2JskCy6g98wEMxpaNNkxohjq3hUKk", fiat: "usd"}, want: 0.0101, want1: 191.64305801999998, wantErr: false},
		{name: "valid eth", args: args{address: "0x76270d9D9afC0cf4EbfFBafE6401E01cb0F021Ce", fiat: "usd"}, want: 0.123, want1: 129.89334508769207, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, _, err := CheckBalance(tt.args.address, tt.args.fiat)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("CheckBalance() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CheckBalance() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
