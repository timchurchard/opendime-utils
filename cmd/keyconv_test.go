package cmd

import (
	"bytes"
	"flag"
	"os"
	"testing"
)

func Test_KeyconvMain(t *testing.T) {
	const (
		cliName                         = "keyconv"
		bitcoinInvalid                  = "Error: WIF malformed/wrong length"
		bitcoinValidCompressedOutput    = "Original WIF: Bitcoin Kx1rJ3afrZvj7jztGupxrtrFoA9SK37CA3ZnwtWDTRt7MQdyvozL compressed=true\n\nBitcoin P2PKH:\t\t\t5HzjtSWNwjML4LiD54MqTFvJQ9R7eAYJwAC4AGJ3NPreEj7B6oN\nBitcoin P2PKH (Compressed):\tKx1rJ3afrZvj7jztGupxrtrFoA9SK37CA3ZnwtWDTRt7MQdyvozL\nBitcoin P2WPKH:\t\t\tp2wpkh:Kx1rJ3afrZvj7jztGupxrtrFoA9SK37CA3ZnwtWDTRt7MQdyvozL\n\nLitecoin P2PKH:\t\t\t6uJUMa3ur9pCXic4at9oEehUMcyaqxzLhqbDsTK55rBFvZ7W4XJ\nLitecoin P2PKH (Compressed):\tT3r7jnsrFwuKtadkpYmq5FPdk1nkP885yFU3oh8m2Q4GsJDSoPhY\nLitecoin P2WPKH:\t\tp2wpkh:T3r7jnsrFwuKtadkpYmq5FPdk1nkP885yFU3oh8m2Q4GsJDSoPhY\n\nDogecoin P2PKH:\t\t\t6JK5BPyLWVe86x9NGSyp4suXsZtf9HpaYLKHUkvwC44H3PzLLJn\n\nEthereum:\t\t\t0x17bc6f773fe98bb1a72adf6e3e89366b7e19ee76d88a023b248fceebfedf1e5d\n"
		bitcoinValidUncompVerboseOutput = "Original WIF: Bitcoin 5JcB6S5ob2WBQdqUC2km8Me9KcDYb9nXGkM5ynL9RascRHfggQs compressed=false\n - Secret exponent: 6831d5d095a1391cebc94315a8f67579ea9c3df8fe8278df682fc536f2f7f907\n\nBitcoin P2PKH:\t\t\t5JcB6S5ob2WBQdqUC2km8Me9KcDYb9nXGkM5ynL9RascRHfggQs\nBitcoin P2PKH (Compressed):\tKziFYyiCdBdnwKaw1gk9mU6wj8SNRjmvcWKeyd5WBWWPKw3EBTqt\nBitcoin P2WPKH:\t\t\tp2wpkh:KziFYyiCdBdnwKaw1gk9mU6wj8SNRjmvcWKeyd5WBWWPKw3EBTqt\n\nLitecoin P2PKH:\t\t\t6uuuZZdLVSy3t1jKhrYiukRKH5n1nxEZ3RkFgyMB93CE75ctC2U\nLitecoin P2PKH (Compressed):\tT6YWzj1P2ZcPiADoZKh1ypeKfz5gVpnpRiDuqRi3kUgYqpbrYbAG\nLitecoin P2WPKH:\t\tp2wpkh:T6YWzj1P2ZcPiADoZKh1ypeKfz5gVpnpRiDuqRi3kUgYqpbrYbAG\n\nDogecoin P2PKH:\t\t\t6JvWPPYm9nnyTFGdPRNjjydNo2h66H4nsvUKJGy3FF5FDwoFpCE\n\nEthereum:\t\t\t0x6831d5d095a1391cebc94315a8f67579ea9c3df8fe8278df682fc536f2f7f907\n"
		bitcoinHexOutput                = "Original WIF: Bitcoin 5KVDhAZchPz2Ywmm3sWvLdPEbKZChuqrpfAehpEy7vmZNosaqgC compressed=false\n\nBitcoin P2PKH:\t\t\t5KVDhAZchPz2Ywmm3sWvLdPEbKZChuqrpfAehpEy7vmZNosaqgC\nBitcoin P2PKH (Compressed):\tL4bZ2HCxZJShqzkZRfy4Rdb8zUu8faEeuMTbR9WehaKfuBwv8QTZ\nBitcoin P2WPKH:\t\t\tp2wpkh:L4bZ2HCxZJShqzkZRfy4Rdb8zUu8faEeuMTbR9WehaKfuBwv8QTZ\n\nLitecoin P2PKH:\t\t\t6vnxAJ79bpSu2KfcZhJt82AQYo7fuiHtbLZpR1FzqP6B4iaR7rR\nLitecoin P2PKH (Compressed):\tTARpU2W8xgRJcqPRyJuvdz8WwLYSjfFYiZMrGx9CGYVqR5Yuj4fE\nLitecoin P2WPKH:\t\tp2wpkh:TARpU2W8xgRJcqPRyJuvdz8WwLYSjfFYiZMrGx9CGYVqR5Yuj4fE\n\nDogecoin P2PKH:\t\t\t6KoYz82aGAGpbZCvFG8txFNU4k2kD388RqHt2JsrwayCBYPGU1z\n\nEthereum:\t\t\t0xdc192045a9261a395445d220890d0969fd7dd2bacec12b9ab3c9827cb0df7bf3\n"
		bitcoinHexOutputAddrs           = "Original WIF: Bitcoin 5KVDhAZchPz2Ywmm3sWvLdPEbKZChuqrpfAehpEy7vmZNosaqgC compressed=false\n\nBitcoin P2PKH:\t\t\t5KVDhAZchPz2Ywmm3sWvLdPEbKZChuqrpfAehpEy7vmZNosaqgC\nBitcoin P2PKH (Compressed):\tL4bZ2HCxZJShqzkZRfy4Rdb8zUu8faEeuMTbR9WehaKfuBwv8QTZ\nBitcoin P2WPKH:\t\t\tp2wpkh:L4bZ2HCxZJShqzkZRfy4Rdb8zUu8faEeuMTbR9WehaKfuBwv8QTZ\n\nLitecoin P2PKH:\t\t\t6vnxAJ79bpSu2KfcZhJt82AQYo7fuiHtbLZpR1FzqP6B4iaR7rR\nLitecoin P2PKH (Compressed):\tTARpU2W8xgRJcqPRyJuvdz8WwLYSjfFYiZMrGx9CGYVqR5Yuj4fE\nLitecoin P2WPKH:\t\tp2wpkh:TARpU2W8xgRJcqPRyJuvdz8WwLYSjfFYiZMrGx9CGYVqR5Yuj4fE\n\nDogecoin P2PKH:\t\t\t6KoYz82aGAGpbZCvFG8txFNU4k2kD388RqHt2JsrwayCBYPGU1z\n\nEthereum:\t\t\t0xdc192045a9261a395445d220890d0969fd7dd2bacec12b9ab3c9827cb0df7bf3\nAddresses for Opendime:\tTODO\n- Bitcoin P2PKH\t\t\t 13dnBMR7AV1CrEZJJFAXe3Q3uiXdXXHGnu \n- Bitcoin P2PKH (Compressed)\t 1A9GW5SxA4qW5ZhymW1uPaj6N4vnVJVWZW \n- Bitcoin P2WPKH\t\t bc1qv3ykvnp4qzktqz65l925jgpe3tpkktvqzfhzrq \n- Ethereum\t\t\t 0xCb19D769c583599DbD7D6D78Eb3279a362672747 \n- Litecoin P2PKH\t\t LMrjSZiwF9FG73FTUP9pv4Tp7vtudfiMG7 \n- Litecoin P2PKH (Compressed)\t LUNDmHknEj5ZLNQ8we1CfbnraHJ4iPRtuG \n- Litecoin P2WPKH\t\t ltc1qv3ykvnp4qzktqz65l925jgpe3tpkktvqx4dxms \n- Dogecoin P2PKH\t\t D7msicMkTtuVPEju2qA6BoZenrFvrMCm24 \n"
	)

	// We manipulate the Args to set them up for the testcases, after this test we restore the initial args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	type args struct {
		flags []string
		key   string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantOut string
	}{
		{"bitcoin valid compressed", args{flags: []string{}, key: "Kx1rJ3afrZvj7jztGupxrtrFoA9SK37CA3ZnwtWDTRt7MQdyvozL"}, 0, bitcoinValidCompressedOutput},
		{"bitcoin invalid", args{flags: []string{}, key: "Kx1rJ3afrZvj7"}, 1, bitcoinInvalid},
		{"bitcoin valid uncomp verbose", args{flags: []string{"-v"}, key: "5JcB6S5ob2WBQdqUC2km8Me9KcDYb9nXGkM5ynL9RascRHfggQs"}, 0, bitcoinValidUncompVerboseOutput},
		{"bitcoin hex", args{flags: []string{}, key: "dc192045a9261a395445d220890d0969fd7dd2bacec12b9ab3c9827cb0df7bf3"}, 0, bitcoinHexOutput},
		{"bitcoin hex addrs", args{flags: []string{"-a"}, key: "dc192045a9261a395445d220890d0969fd7dd2bacec12b9ab3c9827cb0df7bf3"}, 0, bitcoinHexOutputAddrs},
	}
	for _, tt := range tests {
		// reset flags else panic
		flag.CommandLine = flag.NewFlagSet(cliName, flag.ExitOnError)
		os.Args = append([]string{cliName}, tt.args.flags...)

		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}

			if got := KeyconvMain(out, tt.args.key); got != tt.want {
				t.Errorf("KeyconvMain() = %v, want %v", got, tt.want)
			}
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("KeyconvMain() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
