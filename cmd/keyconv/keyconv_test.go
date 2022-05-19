package main

import (
	"bytes"
	"flag"
	"os"
	"testing"
)

func Test_realMain(t *testing.T) {
	const (
		cliName                         = "keyconv"
		bitcoinInvalid                  = "Error: WIF malformed/wrong length"
		bitcoinValidCompressedOutput    = "Original WIF: Bitcoin Kx1rJ3afrZvj7jztGupxrtrFoA9SK37CA3ZnwtWDTRt7MQdyvozL compressed=true\n\nBitcoin P2PKH:\t\t\t5HzjtSWNwjML4LiD54MqTFvJQ9R7eAYJwAC4AGJ3NPreEj7B6oN\nBitcoin P2PKH (Compressed):\tKx1rJ3afrZvj7jztGupxrtrFoA9SK37CA3ZnwtWDTRt7MQdyvozL\nBitcoin P2WPKH:\t\t\tp2wpkh:Kx1rJ3afrZvj7jztGupxrtrFoA9SK37CA3ZnwtWDTRt7MQdyvozL\n\nLitecoin P2PKH:\t\t\t6uJUMa3ur9pCXic4at9oEehUMcyaqxzLhqbDsTK55rBFvZ7W4XJ\nLitecoin P2PKH (Compressed):\tT3r7jnsrFwuKtadkpYmq5FPdk1nkP885yFU3oh8m2Q4GsJDSoPhY\nLitecoin P2WPKH:\t\tp2wpkh:T3r7jnsrFwuKtadkpYmq5FPdk1nkP885yFU3oh8m2Q4GsJDSoPhY\n\nEthereum:\t\t\t0x17bc6f773fe98bb1a72adf6e3e89366b7e19ee76d88a023b248fceebfedf1e5d\n"
		bitcoinValidUncompVerboseOutput = "Original WIF: Bitcoin 5JcB6S5ob2WBQdqUC2km8Me9KcDYb9nXGkM5ynL9RascRHfggQs compressed=false\n - Secret exponent: 6831d5d095a1391cebc94315a8f67579ea9c3df8fe8278df682fc536f2f7f907\n\nBitcoin P2PKH:\t\t\t5JcB6S5ob2WBQdqUC2km8Me9KcDYb9nXGkM5ynL9RascRHfggQs\nBitcoin P2PKH (Compressed):\tKziFYyiCdBdnwKaw1gk9mU6wj8SNRjmvcWKeyd5WBWWPKw3EBTqt\nBitcoin P2WPKH:\t\t\tp2wpkh:KziFYyiCdBdnwKaw1gk9mU6wj8SNRjmvcWKeyd5WBWWPKw3EBTqt\n\nLitecoin P2PKH:\t\t\t6uuuZZdLVSy3t1jKhrYiukRKH5n1nxEZ3RkFgyMB93CE75ctC2U\nLitecoin P2PKH (Compressed):\tT6YWzj1P2ZcPiADoZKh1ypeKfz5gVpnpRiDuqRi3kUgYqpbrYbAG\nLitecoin P2WPKH:\t\tp2wpkh:T6YWzj1P2ZcPiADoZKh1ypeKfz5gVpnpRiDuqRi3kUgYqpbrYbAG\n\nEthereum:\t\t\t0x6831d5d095a1391cebc94315a8f67579ea9c3df8fe8278df682fc536f2f7f907\n"
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
	}
	for _, tt := range tests {
		// reset flags else panic
		flag.CommandLine = flag.NewFlagSet(cliName, flag.ExitOnError)
		os.Args = append([]string{cliName}, tt.args.flags...)

		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}

			if got := realMain(out, tt.args.key); got != tt.want {
				t.Errorf("realMain() = %v, want %v", got, tt.want)
			}
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("realMain() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
