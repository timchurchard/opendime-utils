#!/usr/bin/env python
"""keyconv.py Convert Opendime private key to other formats eg Ethereum and Litecoin

Note: Electrum can import this private key to recover the p2wpkh simply by prefixing it eg p2wpkh:<KEY>
"""

import sys

import click

from pycoin.symbols.btc import network as bitcoin_network
from pycoin.symbols.ltc import network as litecoin_network


@click.command()
@click.option('--key', prompt='Private key WIF', help='Private Key WIF')
def main(key):
    if not key:
        print('Key required')
        sys.exit(1)

    k = bitcoin_network.parse.wif(key)
    print('Bitcoin P2WPKH:\tp2wpkh:%s' % k.as_text())
    print('Ethereum:\t%s' % hex(k.secret_exponent()))
    print('Litecoin:\t%s' % litecoin_network.parse.secret_exponent(k.secret_exponent()).as_text())


if __name__ == '__main__':
    sys.exit(main())
