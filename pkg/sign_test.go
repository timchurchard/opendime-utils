package pkg

import (
	"encoding/hex"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestVerifySignature(t *testing.T) {
	type args struct {
		address   string
		signature []byte
		message   []byte
	}
	tests := []struct {
		name    string
		args    args
		want    VerifiedMessage
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := VerifySignature(tt.args.address, tt.args.signature, tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifySignature() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VerifySignature() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateSignature(t *testing.T) {
	type args struct {
		signature string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateSignature(tt.args.signature)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSignature() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVerifyMessage(t *testing.T) {
	validBitcoinSignatureHex, _ := hex.DecodeString("1f3d290f14b4f5f11d8caad225e9cb62ccf4e606a44fa796f46773283a299e5c142f76dfde641e25c4beb5a90220cb9abcafcd2cea2b5d65d90aa574281325c205")
	validBitcoinMessageHex, _ := hex.DecodeString("48656c6c6f20576f726c64")
	validLitecoinSignatureHex, _ := hex.DecodeString("1f4db5afe1f16d7668d9592ecaad03fe97f3f24965a83ce63b373325705a9c8cad76c6d928f960dead4d872b722d7a74dc341afffd33a0e8e825f5548d9391aa5f")
	validLitecoinMessageHex, _ := hex.DecodeString("48656c6c6f20576f726c64")

	type args struct {
		address   string
		signature string
		message   string
	}
	tests := []struct {
		name    string
		args    args
		want    VerifiedMessage
		wantErr bool
	}{
		{
			name: "valid bitcoin",
			args: args{
				address:   "133r6sCjLq6NbmSLjxuypDuSeocwenu1Qg",
				signature: "Hz0pDxS09fEdjKrSJenLYsz05gakT6eW9GdzKDopnlwUL3bf3mQeJcS+takCIMuavK/NLOorXWXZCqV0KBMlwgU=",
				message:   "Hello World",
			},
			want: VerifiedMessage{
				Address:      "133r6sCjLq6NbmSLjxuypDuSeocwenu1Qg",
				Signature:    validBitcoinSignatureHex,
				Message:      validBitcoinMessageHex,
				IsValid:      true,
				PublicKeyHex: "046afa3afc399f7f332866e37f475938589de0a3298a3aa062c8f4c74450e3d3b287224ee09db6d217912f4706147bb96762d1e11e7ce2e928fd61ecdbd2e37a99",
			},
			wantErr: false,
		}, {
			name: "invalid bitcoin message",
			args: args{
				address:   "133r6sCjLq6NbmSLjxuypDuSeocwenu1Qg",
				signature: "Hz0pDxS09fEdjKrSJenLYsz05gakT6eW9GdzKDopnlwUL3bf3mQeJcS+takCIMuavK/NLOorXWXZCqV0KBMlwgU=",
				message:   "invalid invalid invalid invalid",
			},
			want:    VerifiedMessage{},
			wantErr: true,
		}, {
			name: "invalid bitcoin signature len",
			args: args{
				address:   "133r6sCjLq6NbmSLjxuypDuSeocwenu1Qg",
				signature: "VGhlIHF1aWNrIGJyb3duIGZveCBqdW1wcyBvdmVyIDEzIGxhenkgZG9ncy4=",
				message:   "invalid invalid invalid invalid",
			},
			want:    VerifiedMessage{},
			wantErr: true,
		}, {
			name: "invalid bitcoin signature chars",
			args: args{
				address:   "133r6sCjLq6NbmSLjxuypDuSeocwenu1Qg",
				signature: "Hz0pDxS09fEdjKrSJenLYsz05gakT6eW9GdzKDopnlwUL3bf3mQe????",
				message:   "invalid invalid invalid invalid",
			},
			want:    VerifiedMessage{},
			wantErr: true,
		}, {
			name: "valid litecoin",
			args: args{
				address:   "LLsXEU59RyoMmjgCkUAghxTLr6FXoRCgQT",
				signature: "H021r+HxbXZo2Vkuyq0D/pfz8kllqDzmOzczJXBanIytdsbZKPlg3q1NhytyLXp03DQa//0zoOjoJfVUjZORql8=",
				message:   "Hello World",
			},
			want: VerifiedMessage{
				Address:      "LLsXEU59RyoMmjgCkUAghxTLr6FXoRCgQT",
				Signature:    validLitecoinSignatureHex,
				Message:      validLitecoinMessageHex,
				IsValid:      true,
				PublicKeyHex: "04a2e8f5aa9c46242cdc6463adac2ef8e6bb8b17202c06d17c647066ed143535ac1f93e66cc499170185ec79b2ef5c04119282544fea4c8072ff87711e13597bcf",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := VerifyMessage(tt.args.address, tt.args.signature, tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VerifyMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseVerifyTxt(t *testing.T) {
	const (
		verifyTxtBtc = "-----BEGIN BITCOIN SIGNED MESSAGE-----\nNonce: 1675bf38ec241a2308585ad0  Serial: DDRRNOCZJRIFCIBAEBJDOJQY74\nVersion: 2.4.0 time=20190207.130255 git=master@e233940e coin=BTC\n-----BEGIN SIGNATURE-----\n1Mmg2eycKHomhjAikEAVehHpCSHTREhLfR\nG1pnvdb0RfKfv3Jhg4x0XBQqv1KQx3WFRaxTiUVN84fpIzxOBgapJb/Dpy6auJ28xcHaBxl3XHBbJejfokjgtmg=\n-----END BITCOIN SIGNED MESSAGE-----\n\n\n                                                                    \n                                                                    \n\n"
		verifyTxtLtc = "-----BEGIN LITECOIN SIGNED MESSAGE-----\nUNSEALED -- UNSEALED -- UNSEALED\nNonce: 961f7ecaa917101d4241a43a  Serial: PZZUNUKLGRIFCICKJIYDEEIC74\nVersion: 2.3.0 time=20171018.143523 git=master@8fb7cfd coin=LTC\n-----BEGIN SIGNATURE-----\nLhNxvyyxBGv1Z9CKUaYPE5azvFCMnDMbRN\nHAVOlsYZ4/sj1lVHlqeYd4jbxRRkD5zqp6MG6mNKPmfEdE8rwByiQ+aFTuEpXswhV4y5S5dxREq3pkdq4CjU3/A=\n-----END LITECOIN SIGNED MESSAGE-----\n\n\n                                                                    \n                                   \n"
	)

	btcFp, _ := ioutil.TempFile("", "pvt*")
	btcFp.WriteString(verifyTxtBtc)
	btcFp.Close()
	defer os.Remove(btcFp.Name())

	ltcFp, _ := ioutil.TempFile("", "pvt*")
	ltcFp.WriteString(verifyTxtLtc)
	ltcFp.Close()
	defer os.Remove(ltcFp.Name())

	type args struct {
		fn string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		want2   string
		wantErr bool
	}{
		{
			name: "valid bitcoin",
			args: args{
				fn: btcFp.Name(),
			},
			want:    "1Mmg2eycKHomhjAikEAVehHpCSHTREhLfR",
			want1:   "G1pnvdb0RfKfv3Jhg4x0XBQqv1KQx3WFRaxTiUVN84fpIzxOBgapJb/Dpy6auJ28xcHaBxl3XHBbJejfokjgtmg=",
			want2:   "Nonce: 1675bf38ec241a2308585ad0  Serial: DDRRNOCZJRIFCIBAEBJDOJQY74\r\nVersion: 2.4.0 time=20190207.130255 git=master@e233940e coin=BTC",
			wantErr: false,
		}, {
			name: "valid litecoin",
			args: args{
				fn: ltcFp.Name(),
			},
			want:    "LhNxvyyxBGv1Z9CKUaYPE5azvFCMnDMbRN",
			want1:   "HAVOlsYZ4/sj1lVHlqeYd4jbxRRkD5zqp6MG6mNKPmfEdE8rwByiQ+aFTuEpXswhV4y5S5dxREq3pkdq4CjU3/A=",
			want2:   "UNSEALED -- UNSEALED -- UNSEALED\r\nNonce: 961f7ecaa917101d4241a43a  Serial: PZZUNUKLGRIFCICKJIYDEEIC74\r\nVersion: 2.3.0 time=20171018.143523 git=master@8fb7cfd coin=LTC",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, err := ParseVerifyTxt(tt.args.fn)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseVerifyTxt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseVerifyTxt() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ParseVerifyTxt() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("ParseVerifyTxt() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
