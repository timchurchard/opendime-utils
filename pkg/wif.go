package pkg

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/btcsuite/btcutil/base58"
)

const (
	prefixBitcoin  = 0x80
	prefixLitecoin = 0xb0

	compressedWif = 0x01

	bitcoin  = "Bitcoin"
	litecoin = "Litecoin"
)

// ValidateWif validate WIF is valid and return decoded mode (Bitcoin/Litecoin) secret exponent hex and isCompressed (and error)
func ValidateWif(key string) (string, string, bool, error) {
	var (
		mode           string
		isCompressed   bool
		secretExponent []byte
	)

	const (
		expectedCompressedLen   = 38
		expectedUncompressedLen = 37
	)

	keyBytes := base58.Decode(key)
	keyBytesLen := len(keyBytes)
	switch keyBytesLen {
	case expectedCompressedLen:
		if keyBytes[keyBytesLen-5] != compressedWif {
			return "", "", false, errors.New("WIF malformed/compression byte not 01")
		}
		isCompressed = true
		secretExponent = keyBytes[1 : keyBytesLen-5]
	case expectedUncompressedLen:
		isCompressed = false
		secretExponent = keyBytes[1 : keyBytesLen-4]
	default:
		return "", "", false, errors.New("WIF malformed/wrong length")
	}

	switch keyBytes[0] {
	case prefixBitcoin:
		mode = bitcoin
	case prefixLitecoin:
		mode = litecoin
	default:
		return "", "", false, errors.New("WIF malformed/wrong prefix byte")
	}

	secretExponentHex := hex.EncodeToString(secretExponent)

	return mode, secretExponentHex, isCompressed, nil
}

// ToWif encode a WIF given the prefix, secret exponent hex and compression flag
func ToWif(prefixHex string, secretHex string, compress bool) string {
	bodyHex := prefixHex + secretHex
	if compress {
		bodyHex += "01"
	}

	checksum, err := calcChecksum(bodyHex)
	if err != nil {
		return ""
	}

	wifBytes, _ := hex.DecodeString(bodyHex + checksum)
	return base58.Encode(wifBytes)
}

// calcChecksum takes the wif key hex (prefix + secret exponent + optional is_compressed flag)
// return the checksum hex
func calcChecksum(bodyHex string) (string, error) {
	chkS1, err := hashSha256(bodyHex)
	if err != nil {
		return "", err
	}

	chkS2, _ := hashSha256(chkS1)

	return chkS2[0:8], nil
}

func hashSha256(hexVal string) (string, error) {
	data, err := hex.DecodeString(hexVal)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}
