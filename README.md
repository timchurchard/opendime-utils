# opendime-utils

**Warning! Use at your own risk! I am not responsible for any loss from using this code.**

Utility for making multiple addresses for a sealed [Opendime](https://opendime.com/).

Opendime is a USB stick for storing Bitcoin. A private key is made and stored in a secure chip. The Opendime will not reveal the private key until the end-of-life resistor is removed. Opendimes serve as a bearer instrument letting people safely spend Bitcoin like a coin or note.

Also supports Litecoin Opendimes.

## Premise

It is possible to recover the public key from a Bitcoin/Litecoin signature. And a sealed Opendime can make a signature to prove it controls the private key. Using the public key we are able to derive other addresses such as Bitcoin P2WPKH as well as altcoin addresses such as Litecoin and Ethereum.

## Install

Install the latest from source using the `go install` command

```shell
go install github.com/timchurchard/opendime-utils/...@latest
```

## Usage

This utility provides three sub commands. sigtoaddr to derive addresses from a signature. keyconv to convert a single private key into other formats eg compressed/uncompressed and altcoin formats. crypt to encrypt/decrypt messages using a Bitcoin signature or private key.

## Examples & Tips

sigtoaddr for my tips opendime verify.txt

```shell
$ ./opendime-utils sigtoaddr -verifytxt ./verify.txt_tips
Addresses for Opendime: 1Mmg2eycKHomhjAikEAVehHpCSHTREhLfR
- Bitcoin P2PKH                  1Mmg2eycKHomhjAikEAVehHpCSHTREhLfR
- Bitcoin P2PKH (Compressed)     129azYLPaG55Kb7z1TgvBbj6nRjYFcNMqE
- Bitcoin P2WPKH                 bc1qpjtaggfhsnhkcyg967k3jmsxtm5hzg72q8ejr5
- Ethereum                       0x76270d9D9afC0cf4EbfFBafE6401E01cb0F021Ce
- Litecoin P2PKH                 LfzdHsHSPx3pxXrsvN9nviMaQeejdnT81s
- Litecoin P2PKH (Compressed)    LLNYFkeDevK8aPp9BbgDTcnrze6pQc7D6s
- Litecoin P2WPKH                ltc1qpjtaggfhsnhkcyg967k3jmsxtm5hzg72ymrkmy
```

For an example of crypt see [crypt_demo.sh](crypt_demo.sh).
