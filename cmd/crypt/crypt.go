package main

import (
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	ecies "github.com/ecies/go"

	"github.com/timchurchard/opendime-utils/pkg"
)

func main() {
	os.Exit(realMain(os.Stdout))
}

func realMain(out io.Writer) int { // TODO This func is super ugly, should refactor
	const (
		defaultEmpty    = ""
		usageVerbose    = "Verbose mode"
		usageEncrypt    = "Encrypt"
		usageDecrypt    = "Decrypt"
		usageVerifyTxt  = "Path to OPENDIME/advanced/verify.txt (to encrypt for) alternative to passing (address, signature and message)"
		usageAddress    = "Bitcoin or Litecoin address to encrypt for. Optional with verify.txt"
		usageSignature  = "Bitcoin or Litecoin signature (required if verify.txt not used)"
		usageMessage    = "Bitcoin message (required if verify.txt not used)"
		usageKey        = "Private key in WIF format (for decrypt)"
		usageInput      = "Input string to encrypt/decrypt"
		usageInputFile  = "Path to input file"
		usageOutput     = "Output as string"
		usageOutputFile = "Path to output file"
	)
	var (
		err             error
		verifiedMessage pkg.VerifiedMessage
		verbose         bool
		encrypt         bool
		decrypt         bool
		output          bool
		verifyTxtFn     string
		address         string
		signature       string
		message         string
		privateKey      string
		input           string
		inputFn         string
		outputFn        string
		cipherText      []byte
		plainText       []byte
		rawInput        []byte
	)

	flag.BoolVar(&verbose, "verbose", false, usageVerbose)
	flag.BoolVar(&verbose, "v", false, usageVerbose)

	flag.BoolVar(&encrypt, "encrypt", false, usageEncrypt)
	flag.BoolVar(&encrypt, "e", false, usageEncrypt)

	flag.BoolVar(&decrypt, "decrypt", false, usageDecrypt)
	flag.BoolVar(&decrypt, "d", false, usageDecrypt)

	flag.StringVar(&verifyTxtFn, "verifytxt", defaultEmpty, usageVerifyTxt)

	flag.StringVar(&address, "address", defaultEmpty, usageAddress)
	flag.StringVar(&address, "a", defaultEmpty, usageAddress+" (shorthand)")

	flag.StringVar(&signature, "signature", defaultEmpty, usageSignature)
	flag.StringVar(&signature, "s", defaultEmpty, usageSignature+" (shorthand)")

	flag.StringVar(&message, "message", defaultEmpty, usageMessage)
	flag.StringVar(&message, "m", defaultEmpty, usageMessage+" (shorthand)")

	flag.StringVar(&privateKey, "key", defaultEmpty, usageKey)
	flag.StringVar(&privateKey, "k", defaultEmpty, usageKey+" (shorthand)")

	flag.StringVar(&input, "input", defaultEmpty, usageInput)
	flag.StringVar(&input, "i", defaultEmpty, usageInput+" (shorthand)")

	flag.BoolVar(&output, "output", false, usageOutput)
	flag.BoolVar(&output, "o", false, usageOutput+" (shorthand)")

	flag.StringVar(&inputFn, "inputfile", defaultEmpty, usageInputFile)
	flag.StringVar(&outputFn, "outputfile", defaultEmpty, usageOutputFile)

	flag.Usage = func() {
		fmt.Fprintf(out, "Usage of %s:\n", os.Args[0])

		flag.PrintDefaults()
	}

	flag.Parse()

	// todo terrible sleep because of flag.Parse not finished by the next if ??
	time.Sleep(time.Millisecond * 500)

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
	} else if address == "" && privateKey == "" {
		// verify.txt or (address, signature, message) is required so print usage
		flag.Usage()
		return 1
	}

	if privateKey == "" {
		verifiedMessage, err = pkg.VerifyMessage(address, signature, message)
		if err != nil {
			fmt.Fprintf(out, "Unable to verify signature: %v", err)
			return 1
		}
	}

	if input == "" && inputFn == "" {
		flag.Usage()
		return 1
	}
	if !output && outputFn == "" {
		flag.Usage()
		return 1
	}

	if encrypt && !decrypt {
		publicKey, _ := ecies.NewPublicKeyFromHex(verifiedMessage.PublicKeyHex)

		if input != "" {
			cipherText, err = ecies.Encrypt(publicKey, []byte(input))
		} else {
			data, err := os.ReadFile(inputFn)
			if err != nil {
				fmt.Fprintf(out, "Error reading input file: %s %v", inputFn, err)
				return 1
			}

			cipherText, err = ecies.Encrypt(publicKey, []byte(data))
		}

		if err != nil {
			fmt.Fprintf(out, "Error encrypting data: %v", err)
		}

		fmt.Fprintf(out, "Encrypted message for %s\n", address)
		if output {
			fmt.Fprintf(out, "%s\n", base64.StdEncoding.EncodeToString(cipherText))
		} else {
			err = os.WriteFile(outputFn, cipherText, 0600)
			if err != nil {
				fmt.Fprintf(out, "Error writing output file: %s %v", outputFn, err)
				return 1
			}
			fmt.Fprintf(out, "Written to file: %s\n", outputFn)
		}
	} else if decrypt && !encrypt {
		_, secretHex, _, err := pkg.ValidateWif(privateKey)
		if err != nil {
			fmt.Fprintf(out, "Error decoding WIF: %v", err)
			return 1
		}

		privateKey, err := ecies.NewPrivateKeyFromHex(secretHex)
		if err != nil {
			fmt.Fprintf(out, "Error building private key: %v", err)
			return 1
		}

		if input != "" {
			rawInput, err = base64.StdEncoding.DecodeString(input)
			if err != nil {
				fmt.Fprintf(out, "Error decoding base64 input: %v", err)
				return 1
			}
		} else {
			rawInput, err = os.ReadFile(inputFn)
			if err != nil {
				fmt.Fprintf(out, "Error reading input file: %s %v", inputFn, err)
				return 1
			}
		}

		plainText, err = ecies.Decrypt(privateKey, rawInput)
		if err != nil {
			fmt.Fprintf(out, "Error decrypting: %v", err)
			return 1
		}

		if output {
			fmt.Fprintf(out, "Decrypted message\n%s\n", plainText)
		} else {
			err = os.WriteFile(outputFn, plainText, 0600)
			if err != nil {
				fmt.Fprintf(out, "Error writing output file: %s %v", outputFn, err)
				return 1
			}
			fmt.Fprintf(out, "Decrypted message\nWritten to file: %s\n", outputFn)
		}
	} else {
		flag.Usage()
		return 1
	}

	return 0
}
