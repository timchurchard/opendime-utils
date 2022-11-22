package cmd

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io"

	ecies "github.com/ecies/go"
	"github.com/timchurchard/opendime-utils/pkg"
)

const (
	prefixBitcoinHex  = "80"
	prefixLitecoinHex = "b0"
	prefixDogecoinHex = "9e"
)

// KeyconvMain entrypoint for the keyconv command
func KeyconvMain(out io.Writer, key string) int {
	balance := flag.Bool("b", false, "Show balances")
	makeAddrs := flag.Bool("a", false, "Make addresses")
	verbose := flag.Bool("v", false, "Verbose mode")
	flag.Parse()

	// Check if key is hex and convert to wif
	_, err := hex.DecodeString(key)
	if err == nil {
		key = pkg.ToWif(prefixBitcoinHex, key, false)
	}

	mode, secretExponentHex, isCompressed, err := pkg.ValidateWif(key)
	if err != nil {
		fmt.Fprintf(out, "Error: %v", err)
		return 1
	}

	fmt.Fprintf(out, "Original WIF: %s %s compressed=%v\n", mode, key, isCompressed)
	if *verbose {
		fmt.Fprintf(out, " - Secret exponent: %s\n", secretExponentHex)
	}

	fmt.Fprintln(out, "")
	fmt.Fprintf(out, "Bitcoin P2PKH:\t\t\t%s\n", pkg.ToWif(prefixBitcoinHex, secretExponentHex, false))
	fmt.Fprintf(out, "Bitcoin P2PKH (Compressed):\t%s\n", pkg.ToWif(prefixBitcoinHex, secretExponentHex, true))
	fmt.Fprintf(out, "Bitcoin P2WPKH:\t\t\tp2wpkh:%s\n\n", pkg.ToWif(prefixBitcoinHex, secretExponentHex, true))

	fmt.Fprintf(out, "Litecoin P2PKH:\t\t\t%s\n", pkg.ToWif(prefixLitecoinHex, secretExponentHex, false))
	fmt.Fprintf(out, "Litecoin P2PKH (Compressed):\t%s\n", pkg.ToWif(prefixLitecoinHex, secretExponentHex, true))
	fmt.Fprintf(out, "Litecoin P2WPKH:\t\tp2wpkh:%s\n\n", pkg.ToWif(prefixLitecoinHex, secretExponentHex, true))

	fmt.Fprintf(out, "Dogecoin P2PKH:\t\t\t%s\n\n", pkg.ToWif(prefixDogecoinHex, secretExponentHex, false))

	fmt.Fprintf(out, "Ethereum:\t\t\t0x%s\n", secretExponentHex)

	if *makeAddrs {
		privKey, err := ecies.NewPrivateKeyFromHex(secretExponentHex)
		if err != nil {
			fmt.Fprintf(out, "Failed to make private key: %v", err)
			return 1
		}

		verifiedMessage := pkg.VerifiedMessage{
			Address:      "TODO",
			PublicKeyHex: privKey.PublicKey.Hex(false),
		}

		addresses, err := pkg.GetAddresses(verifiedMessage)
		if err != nil {
			fmt.Fprintf(out, "Failed to make addresses: %v", err)
			return 1
		}

		prettyPrintAddresses(out, addresses, *balance)
	}

	return 0
}
