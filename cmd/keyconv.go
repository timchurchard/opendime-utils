package cmd

import (
	"flag"
	"fmt"
	"io"

	"github.com/timchurchard/opendime-utils/pkg"
)

const (
	prefixBitcoinHex  = "80"
	prefixLitecoinHex = "b0"
)

func KeyconvMain(out io.Writer, key string) int {
	verbose := flag.Bool("v", false, "Verbose mode")
	flag.Parse()

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

	fmt.Fprintf(out, "Ethereum:\t\t\t0x%s\n", secretExponentHex)

	return 0
}
