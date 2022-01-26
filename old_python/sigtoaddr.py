#!/usr/bin/env python
# pylint: disable=invalid-name,line-too-long,missing-function-docstring
"""sigtoaddr: Convert a Bitcoin signature into an address eg Bitcoin Segwit or Altcoin
"""

import base64
import sys
import urllib.request
from base64 import b64decode
from binascii import hexlify
from collections import namedtuple
from hashlib import sha256 as hashlib_sha256
from json import load as json_load
from typing import Optional

import click
from bitcoin.core import Hash160
from bitcoin.core.key import CPubKey
from bitcoin.core.script import CScript, OP_0
from bitcoin.signmessage import BitcoinMessage
from bitcoin.wallet import P2PKHBitcoinAddress, P2WPKHBitcoinAddress
from ecdsa import VerifyingKey, SECP256k1, ellipticcurve, numbertheory
from ecdsa.util import sigdecode_string
from eth_keys import KeyAPI as ETHKeyAPI
from litecoinutils.keys import PublicKey as LitecoinPublicKey
from litecoinutils.keys import add_magic_prefix as litecoin_add_magic_prefix
from litecoinutils.setup import setup as litecoin_setup
# Import pycoin used by the opendime trustme.py to verify the verify.txt. Rather than duplicating the logic.
from pycoin.contrib import msg_signing
from pycoin.ecdsa.secp256k1 import secp256k1_generator
from pycoin.symbols.btc import network as bitcoin_network
from sympy.ntheory import sqrt_mod

# VerifiedMessage: Simple wrapper for a verified bitcoin signature
VerifiedMessage = namedtuple('VerifiedMessage', ['address', 'signature', 'message', 'is_valid', 'public_key'])

# Addresses: Container for addresses derived from the public key
Addresses = namedtuple('Addresses', ['original', 'bitcoin_p2pkh', 'bitcoin_p2pkh_compressed', 'bitcoin_p2wpkh',
                                     'ethereum', 'litecoin_p2pkh', 'litecoin_p2pkh_compressed', 'litecoin_p2wpkh',
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
    if address.startswith('L'):
        return verify_message_litecoin(address, signature, message)

    sig_bytes = base64.b64decode(signature)
    msg_hash = message.GetHash()
    pubkey = CPubKey.recover_compact(msg_hash, sig_bytes)

    return VerifiedMessage(address, signature, message, str(P2PKHBitcoinAddress.from_pubkey(pubkey)) == str(address), pubkey)


def verify_textfile(addr: str, filename: str) -> VerifiedMessage:
    ms = msg_signing.MessageSigner(bitcoin_network, secp256k1_generator)
    with open(filename, encoding='ascii') as fp:
        msg = fp.read().strip()
    msg, sig_addr, sig = ms.parse_signed_message(msg)

    # force newlines to what we need.
    if '\r' not in msg:
        msg = msg.replace('\n', '\r\n')

    if addr and sig_addr != addr:
        raise ValueError('Not signed with correct address')

    if sig_addr.startswith('1'):
        # do math to verify msg
        ok = ms.verify_message(sig_addr, sig, msg)
        if not ok:
            raise ValueError('Invalid or incorrectly-signed verify.txt found.')
        msg = BitcoinMessage(msg)

    return verify_message(sig_addr, sig, msg)


def verify_message_litecoin(address: str, signature: str, message: str) -> VerifiedMessage:  # pylint: disable=too-many-locals
    # Copied from litecoinutils.key PublicKey class verify_message function so we can recover the public key
    """Creates a public key from a message signature and verifies message

    Bitcoin uses a compact format for message signatures (for tx sigs it
    uses normal DER format). The format has the normal r and s parameters
    that ECDSA signatures have but also includes a prefix which encodes
    extra information. Using the prefix the public key can be
    reconstructed from the signature.

    |  Prefix values:
    |      27 - 0x1B = first key with even y
    |      28 - 0x1C = first key with odd y
    |      29 - 0x1D = second key with even y
    |      30 - 0x1E = second key with odd y

    If key is compressed add 4 (31 - 0x1F, 32 - 0x20, 33 - 0x21, 34 - 0x22 respectively)

    Raises
    ------
    ValueError
        If signature is invalid
    """
    sig = b64decode(signature.encode('utf-8'))
    if len(sig) != 65:
        raise ValueError('Invalid signature size')

    # get signature prefix, compressed and recid (which key is odd/even)
    prefix = sig[0]
    if prefix < 27 or prefix > 35:
        return False
    if prefix >= 31:
        compressed = True
        recid = prefix - 31
    else:
        compressed = False
        recid = prefix - 27

    # create message digest -- note double hashing
    message_magic = litecoin_add_magic_prefix(message)
    message_digest = hashlib_sha256( hashlib_sha256(message_magic).digest() ).digest()

    #
    # use recid, r and s to get the point in the curve
    #

    # ECDSA curve using secp256k1 is defined by: y**2 = x**3 + 7
    # This is done modulo p which (secp256k1) is:
    # p is the finite field prime number and is equal to:
    # 2^256 - 2^32 - 2^9 - 2^8 - 2^7 - 2^6 - 2^4 - 1
    # Note that we could also get that from ecdsa lib from the curve, e.g.:
    # SECP256k1.__dict__['curve'].__dict__['_CurveFp__p']
    _p = 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F
    # Curve's a and b are (y**2 = x**3 + a*x + b)
    _a = 0x0000000000000000000000000000000000000000000000000000000000000000
    _b = 0x0000000000000000000000000000000000000000000000000000000000000007
    # Curve's generator point is:
    _Gx = 0x79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798
    _Gy = 0x483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8
    # prime number of points in the group (the order)
    _order = 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141

    # The ECDSA curve (secp256k1) is:
    # Note that we could get that from ecdsa lib, e.g.:
    # SECP256k1.__dict__['curve']
    _curve = ellipticcurve.CurveFp( _p, _a, _b )

    # The generator base point is:
    # Note that we could get that from ecdsa lib, e.g.:
    # SECP256k1.__dict__['generator']
    _G = ellipticcurve.Point( _curve, _Gx, _Gy, _order )

    # get signature's r and s
    r,s = sigdecode_string(sig[1:], _order)

    # ger R's x coordinate
    x = r + (recid // 2) * _order

    # get R's y coordinate (y**2 = x**3 + 7)
    y_values = sqrt_mod( (x**3 + 7) % _p, _p, True )
    if (y_values[0] - recid) % 2 == 0:
        y = y_values[0]
    else:
        y = y_values[1]

    # get R (recovered ephemeral key) from x,y
    R = ellipticcurve.Point(_curve, x, y, _order)

    # get e (hash of message encoded as big integer)
    e = int(hexlify(message_digest), 16)

    # compute public key Q = r^-1 (sR - eG)
    # because Point substraction is not defined we will instead use:
    # Q = r^-1 (sR + (-eG) )
    minus_e = -e % _order
    inv_r = numbertheory.inverse_mod(r, _order)
    Q = inv_r * ( s*R + minus_e*_G )

    # instantiate the public key and verify message
    public_key = VerifyingKey.from_public_point( Q, curve = SECP256k1 )
    key_hex = hexlify(public_key.to_string()).decode('utf-8')
    pubkey = LitecoinPublicKey.from_hex('04' + key_hex)

    valid = True
    if not pubkey.verify(signature, message):
        valid = False

    # confirm that the address provided corresponds to that public key
    if pubkey.get_address(compressed=compressed).to_string() != address:
        valid = False

    return VerifiedMessage(address, signature, message, valid, CPubKey.fromhex(compress_public_key('04' + key_hex)))


def public_key_to_bitcoin_p2pkh(public_key: CPubKey, compressed=False) -> str:
    if compressed:
        public_key = CPubKey.fromhex(compress_public_key(public_key.hex()))
    else:
        public_key = CPubKey.fromhex(decompress_public_key(public_key.hex()))
    return str(P2PKHBitcoinAddress.from_pubkey(public_key))


def public_key_to_bitcoin_p2wpkh(public_key: CPubKey) -> str:
    compressed_public_key = CPubKey.fromhex(compress_public_key(public_key.hex()))
    scriptPubKey = CScript([OP_0, Hash160(compressed_public_key)])
    return str(P2WPKHBitcoinAddress.from_scriptPubKey(scriptPubKey))


def public_key_to_ethereum(public_key: CPubKey) -> str:
    uncompressed_public_key = CPubKey.fromhex(decompress_public_key(public_key.hex()))
    eth_pk = ETHKeyAPI.PublicKey(bytes.fromhex(uncompressed_public_key.hex()[2:]))
    return eth_pk.to_checksum_address()


def public_key_to_litecoin(public_key: CPubKey, compressed=False) -> str:
    litecoin_setup('mainnet')
    compressed_public_key = CPubKey.fromhex(compress_public_key(public_key.hex()))
    litecoin_public_key = LitecoinPublicKey(compressed_public_key.hex())
    return litecoin_public_key.get_address(compressed=compressed).to_string()


def public_key_to_litecoin_p2wpkh(public_key: CPubKey) -> str:
    litecoin_setup('mainnet')
    compressed_public_key = CPubKey.fromhex(compress_public_key(public_key.hex()))
    litecoin_public_key = LitecoinPublicKey(compressed_public_key.hex())
    return litecoin_public_key.get_segwit_address().to_string()


def get_addresses(verified_message: VerifiedMessage) -> Addresses:
    return Addresses(
        original=verified_message.address,
        bitcoin_p2pkh=public_key_to_bitcoin_p2pkh(verified_message.public_key, compressed=False),
        bitcoin_p2pkh_compressed=public_key_to_bitcoin_p2pkh(verified_message.public_key, compressed=True),
        bitcoin_p2wpkh=public_key_to_bitcoin_p2wpkh(verified_message.public_key),
        ethereum=public_key_to_ethereum(verified_message.public_key),
        litecoin_p2pkh=public_key_to_litecoin(verified_message.public_key, compressed=False),
        litecoin_p2pkh_compressed=public_key_to_litecoin(verified_message.public_key, compressed=True),
        litecoin_p2wpkh=public_key_to_litecoin_p2wpkh(verified_message.public_key),
        uncompressed=decompress_public_key(verified_message.public_key.hex()),
        compressed=compress_public_key(verified_message.public_key.hex()),
    )


def check_balance(address: str, fiat: str = 'usd') -> (Optional[float], Optional[float]):
    """check_balance: Basic usage of bitlaps API
    returns (float, float) eg Bitcoin amount (whole Bitcoin not Sats) and Fiat amount
    """
    # Using bitlaps https://developer.bitaps.com/blockchain which is free for 15 reqs in 5s currently
    currency = 'btc'  # if address[0] in ('1', '3', 'b'):
    if address[0].lower() == 'l':
        currency = 'ltc'
    if address[0] == '0':
        currency = 'eth'

    url = f'https://api.bitaps.com/{currency}/v1/blockchain/address/state/{address}'

    data = None
    with urllib.request.urlopen(url) as wp:
        try:
            data = json_load(wp)
        except:
            print('Error reading bitlaps API')
            return (None, None)

    balance = 0
    try:
        balance = data['data']['balance']
    except AttributeError:
        print('Error balance not found')

    if balance:
        if currency == 'eth':
            balance = balance * 0.000000000000000001
        else:
            # Bitcoin & Litecoin has 8 decimal places
            balance = balance * 0.00000001

    price_url = f'https://api.bitaps.com/market/v1/ticker/{currency}{fiat}'

    data = None
    with urllib.request.urlopen(price_url) as wp:
        try:
            data = json_load(wp)
        except:
            print('Error reading bitlaps API')
            return (None, None)

    fiat_value = None
    try:
        fiat_value = balance * float(data['data']['last'])
    except AttributeError:
        print('Error last price not found')

    return (balance, fiat_value)


def pretty_print_addresses(addresses: Addresses, show_balance: bool = False):
    friendly_name = {
        'bitcoin_p2pkh': 'Bitcoin P2PKH\t\t\t',
        'bitcoin_p2pkh_compressed': 'Bitcoin P2PKH (Compressed)\t',
        'bitcoin_p2wpkh': 'Bitcoin P2WPKH\t\t',
        'ethereum': 'Ethereum\t\t\t',
        'litecoin_p2pkh': 'Litecoin P2PKH\t\t',
        'litecoin_p2pkh_compressed': 'Litecoin P2PKH (Compressed)\t',
        'litecoin_p2wpkh': 'Litecoin P2WPKH\t\t',
    }

    print('Addresses for Opendime:\t', addresses.original)
    for name, value in addresses._asdict().items():
        if name not in friendly_name:
            continue

        pad = balance = spacer = fiat = ''
        if show_balance:
            pad = '\t\t'
            if name in ('bitcoin_p2wpkh', 'ethereum', 'litecoin_p2wpkh'):
                pad = '\t'

            spacer = ' = $'
            balance, fiat = check_balance(addresses.ethereum)

        print(f'- {friendly_name[name]} {value}{pad}{balance}{spacer}{fiat}')


def test_get_addresses(address: str, verifytxt: str, signature: str, message: Optional[BitcoinMessage], expected: Addresses):
    if verifytxt:
        verified_message = verify_textfile(address, verifytxt)
    else:
        verified_message = verify_message(address, signature, message)
    assert verified_message.is_valid

    addresses = get_addresses(verified_message)

    assert verified_message.public_key.hex() in (addresses.compressed, addresses.uncompressed)
    assert addresses == expected


def tests():
    test_get_addresses('1Nu1QpfegiGmqHS6YZxkaiGpnqAUXvZz2f', None,
        'HwPlEOxTxs62ruMHZvamv0wmUlbbaY/2ZSqw9Hpdw+FWfgXuSxQ9x55ceSiFyvnlpiZjt+KIhSYnhGnCv8iDe5o=', BitcoinMessage('Hello World!'),
        Addresses(
            original='1Nu1QpfegiGmqHS6YZxkaiGpnqAUXvZz2f',
            bitcoin_p2pkh='1GLfgL9yKVTRRG1D4fdKkEuEQqAE7ob1eB',
            bitcoin_p2pkh_compressed='1Nu1QpfegiGmqHS6YZxkaiGpnqAUXvZz2f',
            bitcoin_p2wpkh='bc1q7qcf63rtp20dsalcwmceucxs0kwn75l95nsxjf',
            ethereum='0x5D0a9F69035Be4275204f9eBbd5cC049e42429c6',
            litecoin_p2pkh='LaZcwYToQ9hUg4hNEocd2Fxzd3XWEMFnQ5',
            litecoin_p2pkh_compressed='Lh7xg2yUmNWq668Fihx3rjLb13XkbHuBMQ',
            litecoin_p2wpkh='ltc1q7qcf63rtp20dsalcwmceucxs0kwn75l9s02z2e',
            uncompressed='0471bb3ef523055565dd5f9864047b9fe93efa10151ff4bb3640f7de6dfdd76cea9d5cb2da17d725a835f25971818e54acc1db69e4866ea23c9dc33f57cb286315',
            compressed='0371bb3ef523055565dd5f9864047b9fe93efa10151ff4bb3640f7de6dfdd76cea',
        )
    )

    test_get_addresses('1Mmg2eycKHomhjAikEAVehHpCSHTREhLfR',
        './verify.txt_tips', None, None,
        Addresses(
            original='1Mmg2eycKHomhjAikEAVehHpCSHTREhLfR',
            bitcoin_p2pkh='1Mmg2eycKHomhjAikEAVehHpCSHTREhLfR',
            bitcoin_p2pkh_compressed='129azYLPaG55Kb7z1TgvBbj6nRjYFcNMqE',
            bitcoin_p2wpkh='bc1qpjtaggfhsnhkcyg967k3jmsxtm5hzg72q8ejr5',
            ethereum='0x76270d9D9afC0cf4EbfFBafE6401E01cb0F021Ce',
            litecoin_p2pkh='LfzdHsHSPx3pxXrsvN9nviMaQeejdnT81s',
            litecoin_p2pkh_compressed='LLNYFkeDevK8aPp9BbgDTcnrze6pQc7D6s',
            litecoin_p2wpkh='ltc1qpjtaggfhsnhkcyg967k3jmsxtm5hzg72ymrkmy',
            uncompressed='04f27deec87586e475f828cb3cd34d2a02a674c204875e91b90ce4ce1e8773289587979932eef0c5f76c5d5fc692db94749e4efba67b692f564190c4b36ca8763a',
            compressed='02f27deec87586e475f828cb3cd34d2a02a674c204875e91b90ce4ce1e87732895',
        )
    )

    test_get_addresses('LLsXEU59RyoMmjgCkUAghxTLr6FXoRCgQT', None,
        'H021r+HxbXZo2Vkuyq0D/pfz8kllqDzmOzczJXBanIytdsbZKPlg3q1NhytyLXp03DQa//0zoOjoJfVUjZORql8=', 'Hello World',
        Addresses(
            original='LLsXEU59RyoMmjgCkUAghxTLr6FXoRCgQT',
            bitcoin_p2pkh='1FZ33nWeZFk2qv8PnCL2VR2wA3KnGchNbZ',
            bitcoin_p2pkh_compressed='12eZyFmKMKZJWvz3aLBPRwPadstFaFGKAF',
            bitcoin_p2wpkh='bc1qzgfsnjuz7972nd9jtqh26qc00ltjns3tjdewkt',
            ethereum='0x33a5f5ff5d6Aeb3152d223C5407C1e71Bb202C76',
            litecoin_p2pkh='LZmzJzpUduz66ipYxLKKmS6hNFh4MgNPKy',
            litecoin_p2pkh_compressed='LLsXEU59RyoMmjgCkUAghxTLr6FXoRCgQT',
            litecoin_p2wpkh='ltc1qzgfsnjuz7972nd9jtqh26qc00ltjns3tk3r2wm',
            uncompressed='04a2e8f5aa9c46242cdc6463adac2ef8e6bb8b17202c06d17c647066ed143535ac1f93e66cc499170185ec79b2ef5c04119282544fea4c8072ff87711e13597bcf',
            compressed='03a2e8f5aa9c46242cdc6463adac2ef8e6bb8b17202c06d17c647066ed143535ac',
        )
    )

    test_get_addresses('LhNxvyyxBGv1Z9CKUaYPE5azvFCMnDMbRN',
        './litecoin_verify.txt_tips', None, None,
        Addresses(
            original='LhNxvyyxBGv1Z9CKUaYPE5azvFCMnDMbRN',
            bitcoin_p2pkh='1PA1fmg86cfxJLWAJSZ5x4XEi2q5kDxpBk',
            bitcoin_p2pkh_compressed='17tcs8A77LNzH3QqwdGjdKcVPiB1Ka3c2j',
            bitcoin_p2wpkh='bc1qfwf7s8qrlcjfulqymrrw3mejnwwas9y5wz5v8r',
            ethereum='0xDdb5Fc6f27921669FCd177f6877A69356dAe889C',
            litecoin_p2pkh='LhNxvyyxBGv1Z9CKUaYPE5azvFCMnDMbRN',
            litecoin_p2pkh_compressed='LS7a8LTwBzd3Xr717mG2uLgFbvYHQbbJ64',
            litecoin_p2wpkh='ltc1qfwf7s8qrlcjfulqymrrw3mejnwwas9y527wgln',
            uncompressed='04db8b0bc1bf85c9727d31b97fc7483b2d9bbc85d57f7e2ed8f617c98a96966271a41db637664355f9c490abd73b8e68a62afb1d40913fc1384f9edb2475009b89',
            compressed='03db8b0bc1bf85c9727d31b97fc7483b2d9bbc85d57f7e2ed8f617c98a96966271',
        )
   )


@click.command()
@click.option('--verifytxt', type=click.Path(exists=True, dir_okay=False, readable=True),
              help='Path to OPENDIME/advanced/verify.txt alternative to passing address, signature and message')
@click.option('--address', '-a', help='Bitcoin address. Optional with verify.txt.')
@click.option('--signature', '-s', help='Bitcoin signature (if verify.txt not used)')
@click.option('--message', '-m', help='Bitcoin message (if verify.txt not used)')
@click.option('--balance', '-b', help='Check balance', count=True)
def main(verifytxt, address, signature, message, balance):
    # Run the sanity tests first
    tests()

    if verifytxt:
        try:
            verified_message = verify_textfile(address, verifytxt)
        except ValueError as ex:
            print(f'Unable to verify text file: {ex}')
            sys.exit(1)
    elif address and signature and message:
        verified_message = verify_message(address, signature, BitcoinMessage(message))
    else:
        print('usage: sigtoaddr.py (--verifytxt and optional --address) OR (--address --signature --message)')
        sys.exit(1)

    addresses = get_addresses(verified_message)
    pretty_print_addresses(addresses, balance)


if __name__ == '__main__':
    sys.exit(main())  # pylint: disable=no-value-for-parameter
