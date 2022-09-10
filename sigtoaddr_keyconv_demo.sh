#!/usr/bin/env bash
#
# opendime-utils sigtoaddr and keyconv demo
#

# Get addresses for legacy address starting 133r
./opendime-utils sigtoaddr --address 133r6sCjLq6NbmSLjxuypDuSeocwenu1Qg --message "Hello World" --signature "Hz0pDxS09fEdjKrSJenLYsz05gakT6eW9GdzKDopnlwUL3bf3mQeJcS+takCIMuavK/NLOorXWXZCqV0KBMlwgU="
#Addresses for Opendime: 133r6sCjLq6NbmSLjxuypDuSeocwenu1Qg
#- Bitcoin P2PKH                  19MkFnavAVX9Njwt43a2sWZrVg9G5jLntU
#- Bitcoin P2PKH (Compressed)     133r6sCjLq6NbmSLjxuypDuSeocwenu1Qg
#- Bitcoin P2WPKH                 bc1qzeapyvz7kl7v5vj865rahts2jjcdz0ssyc3wl8
#- Ethereum                       0x148582B4F60139ce2Bc7E25e7551F31c1122B6f4
#- Litecoin P2PKH                 LTahWztkF9mCdYe3EBZL9XdchtWYC8LJYm
#- Litecoin P2PKH (Compressed)    LMGoN5WZRVLRra8Vv6uH6EyCs1zDmVPhZV
#- Litecoin P2WPKH                ltc1qzeapyvz7kl7v5vj865rahts2jjcdz0ssqyt28h
#- Dogecoin P2PKH                 DDVqo3XZTuRRuk8UndZbRGjTNosZNGHQdo


# Note! keyconv requires interactive input of the private key WIF not cli args
./opendime-utils keyconv
#Private Key WIF: L165TWkVszAp4yHkFsVRj8udU6w2UxfvVMk8bs9QZZyzNmwWVprK
#Original WIF: Bitcoin L165TWkVszAp4yHkFsVRj8udU6w2UxfvVMk8bs9QZZyzNmwWVprK compressed=true
#
#Bitcoin P2PKH:                  5Jh7uE5sVmfviECx7YNr6vSyJ1tfQ6pLNNvGmbvZXVKMFVbFgcJ
#Bitcoin P2PKH (Compressed):     L165TWkVszAp4yHkFsVRj8udU6w2UxfvVMk8bs9QZZyzNmwWVprK
#Bitcoin P2WPKH:                 p2wpkh:L165TWkVszAp4yHkFsVRj8udU6w2UxfvVMk8bs9QZZyzNmwWVprK
#
#Litecoin P2PKH:                 6uzrNMdQQC8oBc6odNAotKE9FVT8buGN94KSUnwbEwdxwQ74jyt
#Litecoin P2PKH (Compressed):    T6vLuG3gHN9QqovcoWSHwVT1QxaLZ3gpJZePTfmx8YA9tfXTpsWk
#Litecoin P2WPKH:                p2wpkh:T6vLuG3gHN9QqovcoWSHwVT1QxaLZ3gpJZePTfmx8YA9tfXTpsWk
#
#Dogecoin P2PKH:                 6K1TCBYq4Xxikqe7JvzpiYSCmSNCuE6byZ3W66ZTM9Wz4AShznB
#
#Ethereum:                       0x736c3aa95b5aacbb1d32dd39ee160b6f8b499082844785bc7a676d6fb793414d
