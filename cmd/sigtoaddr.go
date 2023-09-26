package cmd

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"github.com/timchurchard/opendime-utils/internal"
	"github.com/timchurchard/opendime-utils/pkg"
)

const (
	defaultCurrency = "usd"
	defaultSymbol   = "$"
)

// SigtoaddrMain entrypoint for the sigtoaddr command
func SigtoaddrMain(out io.Writer) int {
	const (
		defaultEmpty   = ""
		usageVerbose   = "Verbose mode"
		usageBalance   = "Check balance"
		usageVerifyTxt = "Path to OPENDIME/advanced/verify.txt alternative to passing address, signature and message"
		usageAddress   = "Bitcoin or Litecoin address. Optional with verify.txt"
		usageSignature = "Bitcoin or Litecoin signature (required if verify.txt not used)"
		usageMessage   = "Bitcoin message (required if verify.txt not used)"
	)
	var (
		err             error
		verifyTxtFn     string
		address         string
		signature       string
		message         string
		verbose         bool
		balance         bool
		verifiedMessage pkg.VerifiedMessage
		addresses       pkg.Addresses
	)

	flag.BoolVar(&verbose, "verbose", false, usageVerbose)
	flag.BoolVar(&verbose, "v", false, usageVerbose)

	flag.BoolVar(&balance, "balance", false, usageBalance)
	flag.BoolVar(&balance, "b", false, usageBalance)

	flag.StringVar(&verifyTxtFn, "verifytxt", defaultEmpty, usageVerifyTxt)

	flag.StringVar(&address, "address", defaultEmpty, usageAddress)
	flag.StringVar(&address, "a", defaultEmpty, usageAddress+" (shorthand)")

	flag.StringVar(&signature, "signature", defaultEmpty, usageSignature)
	flag.StringVar(&signature, "s", defaultEmpty, usageSignature+" (shorthand)")

	flag.StringVar(&message, "message", defaultEmpty, usageMessage)
	flag.StringVar(&message, "m", defaultEmpty, usageMessage+" (shorthand)")

	flag.Usage = func() {
		fmt.Fprintf(out, "Usage of %s:\n", os.Args[0])

		flag.PrintDefaults()
	}

	flag.Parse()

	// Run the sanity tests live at execution time. This matches the old python implementation.
	// The tests do not make calls to the balance API
	if sanityTests() != 0 {
		panic(errors.New("Fatal! sanity tests failed! quitting."))
	}

	if verifyTxtFn != "" {
		address, signature, message, err = pkg.ParseVerifyTxt(verifyTxtFn)
		if errors.Is(err, os.ErrNotExist) {
			fmt.Fprintf(out, "File '%s' not found", verifyTxtFn)
			return 1
		}
		if err != nil {
			fmt.Fprintf(out, "Unable to parse verify.txt: %v", err)
			return 1
		}
	} else if address == "" {
		// verify.txt or (address, signature, message) is required so print usage
		flag.Usage()
		return 1
	}

	verifiedMessage, err = pkg.VerifyMessage(address, signature, message)
	if err != nil {
		fmt.Fprintf(out, "Unable to verify signature: %v", err)
		return 1
	}

	if verbose {
		fmt.Fprintf(out, "Public key hex: %s\n", verifiedMessage.PublicKeyHex)
	}

	addresses, err = pkg.GetAddresses(verifiedMessage)
	if err != nil {
		fmt.Fprintf(out, "Failed to make addresses: %v", err)
		return 1
	}

	prettyPrintAddresses(out, addresses, balance)

	return 0
}

func sanityTests() int {
	var (
		err             error
		address         string
		signature       string
		message         string
		verifiedMessage pkg.VerifiedMessage
		addresses       pkg.Addresses
	)

	// Use working directory
	wd, err := os.Getwd()
	if err != nil {
		return 1
	}
	wdPath := filepath.Dir(wd)

	btcVerifyTxtFn, _ := filepath.Abs(filepath.Join(wdPath, "./verify.txt_tips"))
	ltcVerifyTxtFn, _ := filepath.Abs(filepath.Join(wdPath, "./litecoin_verify.txt_tips"))

	tests := []struct {
		verifytxt string
		address   string
		signature string
		message   string
		addresses pkg.Addresses
	}{
		{
			address: "1Nu1QpfegiGmqHS6YZxkaiGpnqAUXvZz2f", signature: "HwPlEOxTxs62ruMHZvamv0wmUlbbaY/2ZSqw9Hpdw+FWfgXuSxQ9x55ceSiFyvnlpiZjt+KIhSYnhGnCv8iDe5o=", message: "Hello World!",
			addresses: pkg.Addresses{
				Original:                "1Nu1QpfegiGmqHS6YZxkaiGpnqAUXvZz2f",
				BitcoinP2PKH:            "1GLfgL9yKVTRRG1D4fdKkEuEQqAE7ob1eB",
				BitcoinP2PKHCompressed:  "1Nu1QpfegiGmqHS6YZxkaiGpnqAUXvZz2f",
				BitcoinP2WPKH:           "bc1q7qcf63rtp20dsalcwmceucxs0kwn75l95nsxjf",
				Ethereum:                "0x5D0a9F69035Be4275204f9eBbd5cC049e42429c6",
				LitecoinP2PKH:           "LaZcwYToQ9hUg4hNEocd2Fxzd3XWEMFnQ5",
				LitecoinP2PKHCompressed: "Lh7xg2yUmNWq668Fihx3rjLb13XkbHuBMQ",
				LitecoinP2WPKH:          "ltc1q7qcf63rtp20dsalcwmceucxs0kwn75l9s02z2e",
				DogecoinP2PKH:           "DLUmDb6ccuMhxGBooFctJ14qHxtXVWYf4P",
				UncompressedHex:         "0471bb3ef523055565dd5f9864047b9fe93efa10151ff4bb3640f7de6dfdd76cea9d5cb2da17d725a835f25971818e54acc1db69e4866ea23c9dc33f57cb286315",
				CompressedHex:           "0371bb3ef523055565dd5f9864047b9fe93efa10151ff4bb3640f7de6dfdd76cea",
			},
		}, {
			address:   "1Mmg2eycKHomhjAikEAVehHpCSHTREhLfR",
			verifytxt: btcVerifyTxtFn,
			addresses: pkg.Addresses{
				Original:                "1Mmg2eycKHomhjAikEAVehHpCSHTREhLfR",
				BitcoinP2PKH:            "1Mmg2eycKHomhjAikEAVehHpCSHTREhLfR",
				BitcoinP2PKHCompressed:  "129azYLPaG55Kb7z1TgvBbj6nRjYFcNMqE",
				BitcoinP2WPKH:           "bc1qpjtaggfhsnhkcyg967k3jmsxtm5hzg72q8ejr5",
				Ethereum:                "0x76270d9D9afC0cf4EbfFBafE6401E01cb0F021Ce",
				LitecoinP2PKH:           "LfzdHsHSPx3pxXrsvN9nviMaQeejdnT81s",
				LitecoinP2PKHCompressed: "LLNYFkeDevK8aPp9BbgDTcnrze6pQc7D6s",
				LitecoinP2WPKH:          "ltc1qpjtaggfhsnhkcyg967k3jmsxtm5hzg72ymrkmy",
				DogecoinP2PKH:           "DRumZuvFchi4EjMKUpA4CTTR5a1kpHqQXH",
				UncompressedHex:         "04f27deec87586e475f828cb3cd34d2a02a674c204875e91b90ce4ce1e8773289587979932eef0c5f76c5d5fc692db94749e4efba67b692f564190c4b36ca8763a",
				CompressedHex:           "02f27deec87586e475f828cb3cd34d2a02a674c204875e91b90ce4ce1e87732895",
			},
		}, {
			address:   "LLsXEU59RyoMmjgCkUAghxTLr6FXoRCgQT",
			signature: "H021r+HxbXZo2Vkuyq0D/pfz8kllqDzmOzczJXBanIytdsbZKPlg3q1NhytyLXp03DQa//0zoOjoJfVUjZORql8=",
			message:   "Hello World",
			addresses: pkg.Addresses{
				Original:                "LLsXEU59RyoMmjgCkUAghxTLr6FXoRCgQT",
				BitcoinP2PKH:            "1FZ33nWeZFk2qv8PnCL2VR2wA3KnGchNbZ",
				BitcoinP2PKHCompressed:  "12eZyFmKMKZJWvz3aLBPRwPadstFaFGKAF",
				BitcoinP2WPKH:           "bc1qzgfsnjuz7972nd9jtqh26qc00ltjns3tjdewkt",
				Ethereum:                "0x33a5f5ff5d6Aeb3152d223C5407C1e71Bb202C76",
				LitecoinP2PKH:           "LZmzJzpUduz66ipYxLKKmS6hNFh4MgNPKy",
				LitecoinP2PKHCompressed: "LLsXEU59RyoMmjgCkUAghxTLr6FXoRCgQT",
				LitecoinP2WPKH:          "ltc1qzgfsnjuz7972nd9jtqh26qc00ltjns3tk3r2wm",
				DogecoinP2PKH:           "DKh8b3THrfeKNvJzWnKb3BCY3B45ZzCzeH",
				UncompressedHex:         "04a2e8f5aa9c46242cdc6463adac2ef8e6bb8b17202c06d17c647066ed143535ac1f93e66cc499170185ec79b2ef5c04119282544fea4c8072ff87711e13597bcf",
				CompressedHex:           "03a2e8f5aa9c46242cdc6463adac2ef8e6bb8b17202c06d17c647066ed143535ac",
			},
		}, {
			address:   "LhNxvyyxBGv1Z9CKUaYPE5azvFCMnDMbRN",
			verifytxt: ltcVerifyTxtFn,
			addresses: pkg.Addresses{
				Original:                "LhNxvyyxBGv1Z9CKUaYPE5azvFCMnDMbRN",
				BitcoinP2PKH:            "1PA1fmg86cfxJLWAJSZ5x4XEi2q5kDxpBk",
				BitcoinP2PKHCompressed:  "17tcs8A77LNzH3QqwdGjdKcVPiB1Ka3c2j",
				BitcoinP2WPKH:           "bc1qfwf7s8qrlcjfulqymrrw3mejnwwas9y5wz5v8r",
				Ethereum:                "0xDdb5Fc6f27921669FCd177f6877A69356dAe889C",
				LitecoinP2PKH:           "LhNxvyyxBGv1Z9CKUaYPE5azvFCMnDMbRN",
				LitecoinP2PKHCompressed: "LS7a8LTwBzd3Xr717mG2uLgFbvYHQbbJ64",
				LitecoinP2WPKH:          "ltc1qfwf7s8qrlcjfulqymrrw3mejnwwas9y527wgln",
				DogecoinP2PKH:           "DTJ7D2cmQ2aEqLgm32YeVpgqbAZP2QHE5i",
				UncompressedHex:         "04db8b0bc1bf85c9727d31b97fc7483b2d9bbc85d57f7e2ed8f617c98a96966271a41db637664355f9c490abd73b8e68a62afb1d40913fc1384f9edb2475009b89",
				CompressedHex:           "03db8b0bc1bf85c9727d31b97fc7483b2d9bbc85d57f7e2ed8f617c98a96966271",
			},
		},
	}
	for _, tt := range tests {
		if tt.verifytxt != "" {
			address, signature, message, err = pkg.ParseVerifyTxt(tt.verifytxt)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					// Skip this test since can't find verify.txt file
					continue
				}
				return 1
			}

			verifiedMessage, err = pkg.VerifyMessage(address, signature, message)
			if err != nil {
				return 1
			}
		} else {
			verifiedMessage, err = pkg.VerifyMessage(tt.address, tt.signature, tt.message)
			if err != nil {
				return 1
			}
		}

		addresses, err = pkg.GetAddresses(verifiedMessage)
		if err != nil || addresses != tt.addresses {
			return 1
		}
	}

	return 0
}

func prettyPrintAddresses(out io.Writer, addresses pkg.Addresses, balance bool) {
	prefixes := map[string]string{
		"BitcoinP2PKH":            "Bitcoin P2PKH\t\t\t",
		"BitcoinP2PKHCompressed":  "Bitcoin P2PKH (Compressed)\t",
		"BitcoinP2WPKH":           "Bitcoin P2WPKH\t\t",
		"Ethereum":                "Ethereum\t\t\t",
		"LitecoinP2PKH":           "Litecoin P2PKH\t\t",
		"LitecoinP2PKHCompressed": "Litecoin P2PKH (Compressed)\t",
		"LitecoinP2WPKH":          "Litecoin P2WPKH\t\t",
		"DogecoinP2PKH":           "Dogecoin P2PKH\t\t",
	}

	fmt.Fprintf(out, "Addresses for Opendime:\t%s\n", addresses.Original)

	v := reflect.ValueOf(addresses)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldName := typeOfS.Field(i).Name
		if _, ok := prefixes[fieldName]; !ok {
			continue
		}

		address := fmt.Sprintf("%s", v.Field(i).Interface())

		fmt.Fprintf(out, "- %s %s ", prefixes[fieldName], address)

		if balance {
			amount, value, extra, err := internal.CheckBalance(address, defaultCurrency)
			if err != nil {
				// skip price/value print
			} else {
				fmt.Fprintf(out, "(%.08f = %s%.02f)%s", amount, defaultSymbol, value, extra)

				// Put a terrible sleep here to reduce hammering on public/free APIs
				time.Sleep(time.Second / 3)
			}
		}
		fmt.Fprint(out, "\n")
	}
}
