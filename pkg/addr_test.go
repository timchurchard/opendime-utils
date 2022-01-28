package pkg

import (
	"reflect"
	"testing"
)

func TestGetAddresses(t *testing.T) {
	verifiedMessage, _ := VerifyMessage(
		"133r6sCjLq6NbmSLjxuypDuSeocwenu1Qg",
		"Hz0pDxS09fEdjKrSJenLYsz05gakT6eW9GdzKDopnlwUL3bf3mQeJcS+takCIMuavK/NLOorXWXZCqV0KBMlwgU=",
		"Hello World",
	)

	type args struct {
		message VerifiedMessage
	}
	tests := []struct {
		name    string
		args    args
		want    Addresses
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				message: verifiedMessage,
			},
			want: Addresses{
				Original:                "133r6sCjLq6NbmSLjxuypDuSeocwenu1Qg",
				BitcoinP2PKH:            "19MkFnavAVX9Njwt43a2sWZrVg9G5jLntU",
				BitcoinP2PKHCompressed:  "133r6sCjLq6NbmSLjxuypDuSeocwenu1Qg",
				BitcoinP2WPKH:           "bc1qzeapyvz7kl7v5vj865rahts2jjcdz0ssyc3wl8",
				Ethereum:                "0x148582B4F60139ce2Bc7E25e7551F31c1122B6f4",
				LitecoinP2PKH:           "LTahWztkF9mCdYe3EBZL9XdchtWYC8LJYm",
				LitecoinP2PKHCompressed: "LMGoN5WZRVLRra8Vv6uH6EyCs1zDmVPhZV",
				LitecoinP2WPKH:          "ltc1qzeapyvz7kl7v5vj865rahts2jjcdz0ssqyt28h",
				UncompressedHex:         "046afa3afc399f7f332866e37f475938589de0a3298a3aa062c8f4c74450e3d3b287224ee09db6d217912f4706147bb96762d1e11e7ce2e928fd61ecdbd2e37a99",
				CompressedHex:           "036afa3afc399f7f332866e37f475938589de0a3298a3aa062c8f4c74450e3d3b2",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAddresses(tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAddresses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAddresses() = %v, want %v", got, tt.want)
			}
		})
	}
}
