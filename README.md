# opendime-utils

**Warning! Use at your own risk! I am not responsible for any loss from using this code.**

Utility for making multiple addresses for a sealed [Opendime](https://opendime.com/).

Opendime is a USB stick for storing Bitcoin. A private key is made and stored in a secure chip. The Opendime will not reveal the private key until the end-of-life resistor is removed. Opendimes serve as a bearer instrument letting people safely spend Bitcoin like a coin or note.

Also supports Litecoin Opendimes.

## Premise

It should be possible to create multiple addresses associated to the single Bitcoin/Litecoin private key stored on the Opendime.

A sealed Opendime can make a signature to prove it controls the private key. The public key can be recovered from a signature. Using the public key we are able to derive other addresses such as Bitcoin P2WPKH as well as altcoin addresses such as Litecoin and Ethereum.

## Usage

First use sigtoaddr to derive other addresses from an Opendime signature or Opendime (OPENDIME/advanced/verify.txt) file.

Use the optional --balance / -b option to show balance using [bitlaps.com](https://bitaps.com/) API.

Once unsealed the keyconv can be used to convert the private key into formats used for Litecoin or Ethereum. Electrum can import the Bitcoin private key using prefix p2wpkh:<KEY>.

The crypt utility can be used to encrypt a message or file for a Bitcoin address. And to decrypt those messages with the private key.

## Tips!

```shell
$ ./cmd/sigtoaddr/sigtoaddr --verifytxt ./verify.txt_tips
Addresses for Opendime: 1Mmg2eycKHomhjAikEAVehHpCSHTREhLfR
- Bitcoin P2PKH                  1Mmg2eycKHomhjAikEAVehHpCSHTREhLfR
- Bitcoin P2PKH (Compressed)     129azYLPaG55Kb7z1TgvBbj6nRjYFcNMqE
- Bitcoin P2WPKH                 bc1qpjtaggfhsnhkcyg967k3jmsxtm5hzg72q8ejr5
- Ethereum                       0x76270d9D9afC0cf4EbfFBafE6401E01cb0F021Ce
- Litecoin P2PKH                 LfzdHsHSPx3pxXrsvN9nviMaQeejdnT81s
- Litecoin P2PKH (Compressed)    LLNYFkeDevK8aPp9BbgDTcnrze6pQc7D6s
- Litecoin P2WPKH                ltc1qpjtaggfhsnhkcyg967k3jmsxtm5hzg72ymrkmy
```
