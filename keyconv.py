#!/usr/bin/env python
# pylint: disable=invalid-name,line-too-long,missing-function-docstring
"""keyconv.py Convert Opendime private key to other formats eg Ethereum and Litecoin

Note: Electrum can import this private key to recover the p2wpkh simply by prefixing it eg p2wpkh:<KEY>
"""

import sys

import click
from hashlib import sha256 as hashlib_sha256

from base58check import b58decode, b58encode


PREFIX_BITCOIN = 0x80
PREFIX_BITCOIN_STR = '80'
PREFIX_LITECOIN = 0xb0
PREFIX_LITECOIN_STR = 'b0'


def sha256(hexval) :
    byte_array = bytearray.fromhex(hexval)
    m = hashlib_sha256()
    m.update(byte_array)
    return m.hexdigest()


def calc_checksum(k_hex: str) -> str:
    """calc_checksum: Take the wif key hex (prefix + secret exponent + optional is_compressed flag)
    return the checksum hex"""
    assert len(k_hex) in (66, 68)

    chksum_sha256_1 = sha256(k_hex)
    chksum_sha256_2 = sha256(chksum_sha256_1)

    return chksum_sha256_2[0:8]


def toWif(prefix_hex: str, secret_hex: str, compress: bool = False) -> str:
    assert len(prefix_hex) == 2
    assert prefix_hex in (PREFIX_BITCOIN_STR, PREFIX_LITECOIN_STR)
    assert len(secret_hex) == 64

    body = prefix_hex + secret_hex + ('01' if compress else '')
    checksum = calc_checksum(body)
    return b58encode(bytes.fromhex(body + checksum)).decode('ascii')


@click.command()
@click.option('--key', prompt='Private key WIF', help='Private Key WIF')
@click.option('--verbose', '-v', help='Verbose mode', count=True)
def main(key, verbose):
    if not key:
        print('Key required')
        sys.exit(1)

    key_bytes = b58decode(key.encode('ascii'))
    if key_bytes[0] == PREFIX_BITCOIN:
        mode = 'Bitcoin'
    elif key_bytes[0] == PREFIX_LITECOIN:
        mode = 'Litecoin'
    else:
        print('Error! Unsupported WIF')
        sys.exit(1)

    key_bytes_hex = key_bytes.hex()
    is_compressed = len(key_bytes) == 38
    if is_compressed:
        assert key_bytes[-5] == 0x01
        secret_exponent = key_bytes_hex[2:-10]
    else:
        secret_exponent = key_bytes_hex[2:-8]

    print(f'Original WIF: ({mode}) {key} compressed={is_compressed}')
    if verbose:
        print('- WIF Decoded hex = ', key_bytes_hex)

    print('')
    print('Bitcoin P2PKH:\t\t\t', toWif(PREFIX_BITCOIN_STR, secret_exponent, compress=False))
    print('Bitcoin P2PKH (Compressed):\t', toWif(PREFIX_BITCOIN_STR, secret_exponent, compress=True))
    print('Bitcoin P2WPKH:\t\t\t p2wpkh:%s' % toWif(PREFIX_BITCOIN_STR, secret_exponent, compress=True))
    print('')
    print('Litecoin P2PKH:\t\t\t', toWif(PREFIX_LITECOIN_STR, secret_exponent, compress=False))
    print('Litecoin P2PKH (Compressed):\t', toWif(PREFIX_LITECOIN_STR, secret_exponent, compress=True))
    print('Litecoin P2WPKH:\t\t p2wpkh:%s' % toWif(PREFIX_LITECOIN_STR, secret_exponent, compress=True))
    print('')
    print('Ethereum:\t\t\t 0x%s' % key_bytes[1:-5].hex())
    print('')


if __name__ == '__main__':
    sys.exit(main())
