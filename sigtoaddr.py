#!/usr/bin/env python
"""sigtoaddr: Convert a Bitcoin signature into an address eg Bitcoin Segwit or Altcoin
"""

import sys
from collections import namedtuple

import base64
import click

from bitcoin.core import Hash160
from bitcoin.core.key import CPubKey
from bitcoin.core.script import CScript, OP_0
from bitcoin.signmessage import BitcoinMessage
from bitcoin.wallet import P2PKHBitcoinAddress, P2WPKHBitcoinAddress

from eth_keys import KeyAPI as ETHKeyAPI

from litecoinutils.setup import setup as litecoin_setup
from litecoinutils.keys import PublicKey as LitecoinPublicKey

# Import pycoin used by the opendime trustme.py to verify the verify.txt. Rather than duplicating the logic.
from pycoin.ecdsa.secp256k1 import secp256k1_generator
from pycoin.symbols.btc import network as bitcoin_network
from pycoin.contrib import msg_signing


# VerifiedMessage: Simple wrapper for a verified bitcoin signature
VerifiedMessage = namedtuple('VerifiedMessage', ['address', 'signature', 'message', 'is_valid', 'public_key'])

# Addresses: Container for addresses derived from the public key
Addresses = namedtuple('Addresses', ['bitcoin_p2pkh', 'bitcoin_p2wpkh', 'ethereum', 'litecoin', 'litecoin_p2wpkh',
                                     'uncompressed', 'compressed'])


def compress_public_key(uncompressed_public_key_hex: str) -> str:
    if uncompressed_public_key_hex[0:2] in ('02', '03'):
        # Already compressed
        return uncompressed_public_key_hex

    xi = int(uncompressed_public_key_hex[2:66], 16)
    x = int.to_bytes(xi, length=32, byteorder='big', signed=False)
    yi = int(uncompressed_public_key_hex[66:], 16)
    header = b'\x03' if yi & 1 else b'\x02'
    return (header + x).hex()


def decompress_public_key(compressed_public_key_hex) -> str:
    if compressed_public_key_hex[0:2] == '04':
        # Already uncompressed
        return compressed_public_key_hex

    p = 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F
    compressed_public_key_bytes = bytes.fromhex(compressed_public_key_hex)
    x = int.from_bytes(compressed_public_key_bytes[1:33], byteorder='big')
    y_sq = (pow(x, 3, p) + 7) % p
    y = pow(y_sq, (p + 1) // 4, p)
    if y % 2 != compressed_public_key_bytes[0] % 2:
        y = p - y
    y = y.to_bytes(32, byteorder='big')
    return (b'\x04' + compressed_public_key_bytes[1:33] + y).hex()


def verify_message(address: str, signature: str, message: BitcoinMessage) -> VerifiedMessage:
    sig_bytes = base64.b64decode(signature)
    msg_hash = message.GetHash()
    pubkey = CPubKey.recover_compact(msg_hash, sig_bytes)

    return VerifiedMessage(address, signature, message, str(P2PKHBitcoinAddress.from_pubkey(pubkey)) == str(address), pubkey)


def verify_textfile(addr: str, filename: str) -> VerifiedMessage:
    ms = msg_signing.MessageSigner(bitcoin_network, secp256k1_generator)
    msg = open(filename).read().strip()
    msg, sig_addr, sig = ms.parse_signed_message(msg)

    # force newlines to what we need.
    if '\r' not in msg:
        msg = msg.replace('\n', '\r\n')

    if sig_addr != addr:
        raise ValueError('Not signed with correct address')

    # do math to verify msg
    ok = ms.verify_message(addr, sig, msg)

    if not ok:
        raise ValueError('Invalid or incorrectly-signed verify.txt found.')

    return verify_message(addr, sig, BitcoinMessage(msg))


def public_key_to_bitcoin_p2wpkh(public_key: CPubKey) -> str:
    compressed_public_key = CPubKey.fromhex(compress_public_key(public_key.hex()))
    scriptPubKey = CScript([OP_0, Hash160(compressed_public_key)])
    return str(P2WPKHBitcoinAddress.from_scriptPubKey(scriptPubKey))


def public_key_to_ethereum(public_key: CPubKey) -> str:
    uncompressed_public_key = CPubKey.fromhex(decompress_public_key(public_key.hex()))
    eth_pk = ETHKeyAPI.PublicKey(bytes.fromhex(uncompressed_public_key.hex()[2:]))
    return eth_pk.to_checksum_address()


def public_key_to_litecoin(public_key: CPubKey) -> str:
    litecoin_setup('mainnet')
    compressed_public_key = CPubKey.fromhex(compress_public_key(public_key.hex()))
    litecoin_public_key = LitecoinPublicKey(compressed_public_key.hex())
    return litecoin_public_key.get_address().to_string()


def public_key_to_litecoin_p2wpkh(public_key: CPubKey) -> str:
    litecoin_setup('mainnet')
    compressed_public_key = CPubKey.fromhex(compress_public_key(public_key.hex()))
    litecoin_public_key = LitecoinPublicKey(compressed_public_key.hex())
    return litecoin_public_key.get_segwit_address().to_string()


def get_addresses(verified_message: VerifiedMessage) -> Addresses:
    # todo: deal with magic order and failures and many dereferences of .public_key
    return Addresses(
        verified_message.address,
        public_key_to_bitcoin_p2wpkh(verified_message.public_key),
        public_key_to_ethereum(verified_message.public_key),
        public_key_to_litecoin(verified_message.public_key),
        public_key_to_litecoin_p2wpkh(verified_message.public_key),
        decompress_public_key(verified_message.public_key.hex()),
        compress_public_key(verified_message.public_key.hex()),
    )


def pretty_print_addresses(addresses: Addresses):
    print('Addresses for Opendime:\t', addresses.bitcoin_p2pkh)
    print('- Bitcoin P2WPKH:\t', addresses.bitcoin_p2wpkh)
    print('- Ethereum:\t\t', addresses.ethereum)
    print('- Litecoin:\t\t', addresses.litecoin)
    print('- Litecoin P2WPKH:\t', addresses.litecoin_p2wpkh)


def test_get_addresses(address, verifytxt, signature, message, p2wpkh_expected, ethereum_expected, litecoin_expected, litecoin_p2wpkh_expected):
    if verifytxt:
        verified_message = verify_textfile(address, verifytxt)
    else:
        verified_message = verify_message(address, signature, message)
    assert verified_message.is_valid

    addresses = get_addresses(verified_message)
    assert verified_message.public_key.hex() in (addresses.compressed, addresses.uncompressed)
    assert addresses.bitcoin_p2pkh == address
    assert addresses.bitcoin_p2wpkh == p2wpkh_expected
    assert addresses.ethereum == ethereum_expected
    assert addresses.litecoin == litecoin_expected
    assert addresses.litecoin_p2wpkh == litecoin_p2wpkh_expected


def tests():
    # todo: more tests
    test_get_addresses(
        '1Nu1QpfegiGmqHS6YZxkaiGpnqAUXvZz2f',
        None,
        'HwPlEOxTxs62ruMHZvamv0wmUlbbaY/2ZSqw9Hpdw+FWfgXuSxQ9x55ceSiFyvnlpiZjt+KIhSYnhGnCv8iDe5o=',
        BitcoinMessage('Hello World!'),
        'bc1q7qcf63rtp20dsalcwmceucxs0kwn75l95nsxjf',
        '0x5D0a9F69035Be4275204f9eBbd5cC049e42429c6',
        'Lh7xg2yUmNWq668Fihx3rjLb13XkbHuBMQ',
        'ltc1q7qcf63rtp20dsalcwmceucxs0kwn75l9s02z2e')

    test_get_addresses(
        '1Mmg2eycKHomhjAikEAVehHpCSHTREhLfR',
        './verify.txt_tips', None, None,
        'bc1qpjtaggfhsnhkcyg967k3jmsxtm5hzg72q8ejr5',
        '0x76270d9D9afC0cf4EbfFBafE6401E01cb0F021Ce',
        'LLNYFkeDevK8aPp9BbgDTcnrze6pQc7D6s',
        'ltc1qpjtaggfhsnhkcyg967k3jmsxtm5hzg72ymrkmy')


@click.command()
@click.option('--verifytxt', type=click.Path(exists=True, dir_okay=False, readable=True),
              help='Path to OPENDIME/advanced/verify.txt alternative to passing address, signature and message')
@click.option('--address', '-a', help='Bitcoin address. Required.')
@click.option('--signature', '-s', help='Bitcoin signature (if verify.txt not used)')
@click.option('--message', '-m', help='Bitcoin message (if verify.txt not used)')
def main(verifytxt, address, signature, message):
    # Run the sanity tests first
    tests()

    verified_message = None
    if verifytxt:
        try:
            verified_message = verify_textfile(address, verifytxt)
        except ValueError as ex:
            print('Unable to verify text file: %s' % ex)
            sys.exit(1)
    elif address and signature and message:
        verified_message = verify_message(address, signature, BitcoinMessage(message))
    else:
        print('(--verifytxt and --address) OR (--address --signature --message) are required.')
        sys.exit(1)

    addresses = get_addresses(verified_message)
    pretty_print_addresses(addresses)


if __name__ == '__main__':
    sys.exit(main())
