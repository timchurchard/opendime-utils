# opendime-utils

![Build Status](https://github.com/timchurchard/opendime-utils/workflows/Test/badge.svg)
![Coverage](https://img.shields.io/badge/Coverage-65.0%25-yellow)
[![License](https://img.shields.io/github/license/timchurchard/opendime-utils)](/LICENSE)
[![Release](https://img.shields.io/github/release/timchurchard/opendime-utils.svg)](https://github.com/timchurchard/opendime-utils/releases/latest)
[![GitHub Releases Stats of opendime-utils](https://img.shields.io/github/downloads/timchurchard/opendime-utils/total.svg?logo=github)](https://somsubhra.github.io/github-release-stats/?username=timchurchard&repository=opendime-utils)

**Warning! Use at your own risk! I am not responsible for any loss from using this code.**

Utility for making multiple addresses for a sealed [Opendime](https://opendime.com/). Or from any Bitcoin/Litecoin signed message.

Opendime is a USB stick for storing Bitcoin. A private key is made and stored in a secure chip. The Opendime will not reveal the private key until the end-of-life resistor is removed. Opendimes serve as a bearer instrument letting people safely spend Bitcoin like a coin or note.

Also supported:
- Litecoin Opendimes.
- Encrypt/Decrypt message or file for opendime or any Bitcoin/Litecoin signed message.

## Premise

It is possible to recover the public key from a Bitcoin/Litecoin signature. And a sealed Opendime can make a signature to prove it controls the private key. Using the public key we are able to derive other addresses such as Bitcoin P2WPKH as well as altcoin addresses such as Litecoin, Ethereum and Dogecoin.

## Install

Install the latest from source using the `go install` command

```shell
go install github.com/timchurchard/opendime-utils@latest
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
- Dogecoin P2PKH                 DRumZuvFchi4EjMKUpA4CTTR5a1kpHqQXH
```

For a full example of sigtoaddr and keyconv [here](sigtoaddr_keyconv_demo.sh)

For an example of crypt see [crypt_demo.sh](crypt_demo.sh)

## Encryption with ECIES

**WARNING** Version 0.1.x used ECIES v1 which is incompatible with ECIES v2 in the 0.2.x releases. If you need to decrypt a v1 encrypted message you will need to use a v1 release.
