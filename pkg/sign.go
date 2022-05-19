package pkg

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"syscall"

	"github.com/btcsuite/btcd/wire"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

// VerifiedMessage struct holds response from ValidateSignature
type VerifiedMessage struct {
	Address      string
	Signature    []byte
	Message      []byte
	IsValid      bool
	PublicKeyHex string
}

const (
	// Expected length of verify.txt file
	expectedVerifyTxtLen = 512
	// Expected signature decoded/bytes length
	expectedSignatureLen = 65

	// verify.txt constants
	vtHeaderBitcoin  = "-----BEGIN BITCOIN SIGNED MESSAGE-----\n"
	vtHeaderLitecoin = "-----BEGIN LITECOIN SIGNED MESSAGE-----\n"
	vtSignedMessage  = "\n-----BEGIN SIGNATURE-----\n"
	vtFooterPrefix   = "\n-----END "
)

// ValidateSignature takes a Bitcoin/Litecoin encoded signature and returns the 65 byte DER encoded bytes
func ValidateSignature(signature string) ([]byte, error) {
	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return nil, err
	}

	if len(signatureBytes) != expectedSignatureLen {
		return nil, fmt.Errorf("signature bytes wrong length expected 65 got %d", len(signatureBytes))
	}

	return signatureBytes, nil
}

// VerifyMessage wrapper for VerifySignature that accepts strings for signature and message
// signature will be in Bitcoin base58 format and message as string
func VerifyMessage(address string, signature string, message string) (VerifiedMessage, error) {
	signatureBytes, err := ValidateSignature(signature)
	if err != nil {
		return VerifiedMessage{}, err
	}

	return VerifySignature(address, signatureBytes, []byte(message))
}

// VerifySignature takes an address, signature and message and returns VerifiedMessage
func VerifySignature(address string, signature []byte, message []byte) (VerifiedMessage, error) {
	var buf bytes.Buffer

	bitcoinSignatureHeader := []byte("Bitcoin Signed Message:\n")
	litecoinSignatureHeader := []byte("Litecoin Signed Message:\n")

	btcMatch, _ := regexp.MatchString("^(1|3|bc1).", address)
	if btcMatch { // todo: not sure if 3 addresses can sign? And litecoin has 3 addresses too?
		_ = wire.WriteVarBytes(&buf, 0, bitcoinSignatureHeader)
	} else { // todo: handle other signature types that Bitcoin/Litecoin?
		_ = wire.WriteVarBytes(&buf, 0, litecoinSignatureHeader)
	}

	_ = wire.WriteVarBytes(&buf, 0, message)
	messageHash := chainhash.DoubleHashB(buf.Bytes())

	publicKey, _, err := btcec.RecoverCompact(btcec.S256(), signature, messageHash)
	if err != nil {
		return VerifiedMessage{}, err
	}

	addrs, _ := GetAddresses(VerifiedMessage{PublicKeyHex: hex.EncodeToString(publicKey.SerializeUncompressed())})
	if address != addrs.BitcoinP2PKH && address != addrs.BitcoinP2PKHCompressed &&
		address != addrs.LitecoinP2PKH && address != addrs.LitecoinP2PKHCompressed {
		return VerifiedMessage{}, errors.New("Invalid signature address not match")
	}

	return VerifiedMessage{
		Address:      address,
		Signature:    signature,
		Message:      message,
		IsValid:      true,
		PublicKeyHex: hex.EncodeToString(publicKey.SerializeUncompressed()),
	}, nil
}

// ParseVerifyTxt takes the opendime verify.txt format and returns address, signature, message
// based on https://github.com/richardkiss/pycoin/blob/main/pycoin/contrib/msg_signing.py
func ParseVerifyTxt(fn string) (string, string, string, error) {
	var (
		err       error
		address   string
		message   string
		signature string
	)

	file, err := os.OpenFile(fn, os.O_RDONLY|syscall.O_NOATIME, 0)
	if err != nil {
		return "", "", "", err
	}

	bufb := make([]byte, expectedVerifyTxtLen)
	_, err = file.Read(bufb)
	if err != nil {
		return "", "", "", err
	}

	bufs := string(bufb)

	// Strip any windows line endings to make the string easier to work with
	bufs = strings.ReplaceAll(bufs, "\r\n", "\n")

	// Get the message
	if strings.Index(bufs, vtHeaderBitcoin) == 0 {
		message = bufs[len(vtHeaderBitcoin):strings.Index(bufs, vtSignedMessage)]
	} else if strings.Index(bufs, vtHeaderLitecoin) == 0 {
		message = bufs[len(vtHeaderLitecoin):strings.Index(bufs, vtSignedMessage)]
	} else {
		return "", "", "", fmt.Errorf("verify.txt does not start with '%s' or '%s'", vtHeaderBitcoin, vtHeaderLitecoin)
	}

	// Ensure the final message has DOS newlines as specified in the RFC
	message = strings.ReplaceAll(message, "\n", "\r\n")

	// Get the address and signature
	sigAreaStart := strings.Index(bufs, vtSignedMessage) + len(vtSignedMessage)
	sigAreaEnd := sigAreaStart + strings.Index(bufs[sigAreaStart:], vtFooterPrefix)
	sigArea := bufs[sigAreaStart:sigAreaEnd]

	sigAreaSplit := strings.Split(sigArea, "\n")
	if len(sigAreaSplit) < 2 {
		return "", "", "", fmt.Errorf("verify.txt does not end with address and signature")
	}

	address = sigAreaSplit[0]
	signature = sigAreaSplit[1]

	return address, signature, message, nil
}
