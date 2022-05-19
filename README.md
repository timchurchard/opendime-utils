# opendime-utils

Utilities for making multiple addresses for a sealed [opendime](https://opendime.com/).

Opendime is a small USB stick for storing Bitcoin. A private key is made and stored in a secure chip. The opendime will not reveal the private key until the end-of-life resistor is removed. Opendimes serve as a bearer instrument letting people safely spend Bitcoin like a coin or note.

## Premise

It should be possible to create multiple addresses associated to the single Bitcoin private key stored on the opendime.

A sealed opendime can make a signature to prove it controls the private key. The public key can be recovered from the signature. Using the public key we are able to derive other addresses such as Bitcoin P2WPKH as well as altcoin addresses such as Litecoin and Ethereum.

## Usage

Create a python venv and install the requirements.

```bash
python3 -m venv venv
. venv/bin/activate
python3 -m pip install -U pip setuptools
pip3 install -r requirements.txt
```

First use sigtoaddr.py to derive other addresses from an opendime signature or opendime (OPENDIME/advanced/verify.txt) file. The following example shows creating address for my 1Mmg2 Opendime using its verify.txt. Note the optional --balance -b option to show balance.

```bash
python3 sigtoaddr.py --address 1Mmg2eycKHomhjAikEAVehHpCSHTREhLfR --verifytxt ./verify.txt_tips
```

Once unsealed the keyconv.py can be used to convert the private key into formats used for Litecoin or Ethereum. Electrum can import the Bitcoin private key using prefix p2wpkh:<KEY>.

## Tips!

Addresses for Opendime:  1Mmg2eycKHomhjAikEAVehHpCSHTREhLfR
- Bitcoin P2WPKH:        bc1qpjtaggfhsnhkcyg967k3jmsxtm5hzg72q8ejr5
- Ethereum:              0x76270d9D9afC0cf4EbfFBafE6401E01cb0F021Ce
- Litecoin:              LLNYFkeDevK8aPp9BbgDTcnrze6pQc7D6s
- Litecoin P2WPKH:       ltc1qpjtaggfhsnhkcyg967k3jmsxtm5hzg72ymrkmy
