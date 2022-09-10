package pkg

import (
	"encoding/hex"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// Addresses struct for GetAddresses response
type Addresses struct {
	Original                string
	BitcoinP2PKH            string
	BitcoinP2PKHCompressed  string
	BitcoinP2WPKH           string
	Ethereum                string
	LitecoinP2PKH           string
	LitecoinP2PKHCompressed string
	LitecoinP2WPKH          string
	DogecoinP2PKH           string
	UncompressedHex         string
	CompressedHex           string
}

// GetAddresses get addresses from a verified message (only public key is needed)
func GetAddresses(message VerifiedMessage) (Addresses, error) {
	publicKeyBytes, err := hex.DecodeString(message.PublicKeyHex)
	if err != nil {
		return Addresses{}, err
	}
	publicKey, err := btcec.ParsePubKey(publicKeyBytes, btcec.S256())
	if err != nil {
		return Addresses{}, err
	}

	// todo: assuming these cannot error for a valid public key on the Secp256k1 curve ?!
	pkHash := btcutil.Hash160(publicKey.SerializeUncompressed())
	bitcoinP2PKH, _ := btcutil.NewAddressPubKeyHash(pkHash, &chaincfg.MainNetParams)

	pkHashC := btcutil.Hash160(publicKey.SerializeCompressed())
	bitcoinP2PKHC, _ := btcutil.NewAddressPubKeyHash(pkHashC, &chaincfg.MainNetParams)
	bitcoinP2WPKH, _ := btcutil.NewAddressWitnessPubKeyHash(pkHashC, &chaincfg.MainNetParams)

	ltcPkHash := btcutil.Hash160(publicKey.SerializeUncompressed())
	litecoinP2PKH, _ := btcutil.NewAddressPubKeyHash(ltcPkHash, &LitecoinMainNetParams)

	ltcPkHashC := btcutil.Hash160(publicKey.SerializeCompressed())
	litecoinP2PKHC, _ := btcutil.NewAddressPubKeyHash(ltcPkHashC, &LitecoinMainNetParams)
	litecoinP2WPKH, _ := btcutil.NewAddressWitnessPubKeyHash(ltcPkHashC, &LitecoinMainNetParams)

	dogecoinP2PKH, _ := btcutil.NewAddressPubKeyHash(pkHash, &DogecoinMainNetParams)

	return Addresses{
		Original:                message.Address,
		BitcoinP2PKH:            bitcoinP2PKH.String(),
		BitcoinP2PKHCompressed:  bitcoinP2PKHC.String(),
		BitcoinP2WPKH:           bitcoinP2WPKH.String(),
		Ethereum:                crypto.PubkeyToAddress(*publicKey.ToECDSA()).Hex(),
		LitecoinP2PKH:           litecoinP2PKH.String(),
		LitecoinP2PKHCompressed: litecoinP2PKHC.String(),
		LitecoinP2WPKH:          litecoinP2WPKH.String(),
		DogecoinP2PKH:           dogecoinP2PKH.String(),
		UncompressedHex:         hex.EncodeToString(publicKey.SerializeUncompressed()),
		CompressedHex:           hex.EncodeToString(publicKey.SerializeCompressed()),
	}, nil
}
