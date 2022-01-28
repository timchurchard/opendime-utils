package pkg

import (
	"testing"
)

func Test_ValidateWif(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		want2   bool
		wantErr bool
	}{
		{"invalid wrong compression byte", args{key: "KwUhWnQRC7mKrEvMmYjzs2Qtw3LLktTQjL8GmjNbvVi8vjxYZZAd"}, "", "", false, true},
		{"invalid wrong prefix byte", args{key: "4XAZUtEJi962ESjyenv7vUU7qwzQ1WjJDVxLv5d7wVBNQ6CWD5e"}, "", "", false, true},
		{"invalid too short", args{key: "KwUhWnQRC7mK"}, "", "", false, true},
		{"invalid too long", args{key: "KwUhWnQRC7mKrEvMmYjzs2Qtw3LLktTQjL8GmjNbvVi8vM677LgMKKKKK"}, "", "", false, true},
		{
			"valid bitcoin compressed fc3f",
			args{key: "L5g3omnu8BYUS5zUA74AW1eSbZ1xx72HzSVgJcejsvMTn3P579qd"},
			"Bitcoin", "fc3fa47324ceb77e1160833eddd30ea15efa22a6e59c204921e12fbbab1becb8", true, false,
		},
		{
			"valid bitcoin uncompressed fc3f",
			args{key: "5KjNw6cmtUK1KpoYytfnCZKTC11DgDhAjvMYZYBpKncuHd6YzkX"},
			"Bitcoin", "fc3fa47324ceb77e1160833eddd30ea15efa22a6e59c204921e12fbbab1becb8", false, false,
		},
		{
			"valid litecoin compressed 07b5",
			args{key: "T3JxxXhbbVjvd5ZEKBgs5NxGstyepyUJYY2XdY19VTtJSEep3SM3"},
			"Litecoin", "07b5eb6760c9b0cef7009acc4b2f01d847a5da1e7aa97373f7db996db295ed26", true, false,
		},
		{
			"valid litecoin uncompressed 07b5",
			args{key: "6uBR1M6aDM76oh141wz4eF36s3iPAWZU5syzEiTmX1euSDQkoLG"},
			"Litecoin", "07b5eb6760c9b0cef7009acc4b2f01d847a5da1e7aa97373f7db996db295ed26", false, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, err := ValidateWif(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateWif() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateWif() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ValidateWif() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("ValidateWif() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func Test_ToWif(t *testing.T) {
	type args struct {
		prefixHex string
		secretHex string
		compress  bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"valid bitcoin uncompressed", args{prefixHex: "80", secretHex: "07b5eb6760c9b0cef7009acc4b2f01d847a5da1e7aa97373f7db996db295ed26", compress: false}, "5HsgYDZ3JveELK7CW8C6rrFvua9uxi7SKCapXXSjoZLHkTCgnGD"},
		{"valid bitcoin compressed", args{prefixHex: "80", secretHex: "07b5eb6760c9b0cef7009acc4b2f01d847a5da1e7aa97373f7db996db295ed26", compress: true}, "KwUhWnQRC7mKrEvMmYjzs2Qtw3LLktTQjL8GmjNbvVi8vM677LgM"},
		{"invalid bitcoin", args{prefixHex: "80", secretHex: "abcdefghijkl", compress: false}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToWif(tt.args.prefixHex, tt.args.secretHex, tt.args.compress); got != tt.want {
				t.Errorf("ToWif() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calcChecksum(t *testing.T) {
	type args struct {
		bodyHex string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"invalid", args{bodyHex: "abcdefghijkl"}, "", true},
		{"valid", args{bodyHex: "deadbeef"}, "281dd50f", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calcChecksum(tt.args.bodyHex)
			if (err != nil) != tt.wantErr {
				t.Errorf("calcChecksum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("calcChecksum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hashSha256(t *testing.T) {
	type args struct {
		hexVal string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"invalid", args{hexVal: "abcdefghijkl"}, "", true},
		{"valid", args{hexVal: "deadbeef"}, "5f78c33274e43fa9de5659265c1d917e25c03722dcb0b8d27db8d5feaa813953", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := hashSha256(tt.args.hexVal)
			if (err != nil) != tt.wantErr {
				t.Errorf("hashSha256() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("hashSha256() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateWif(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		want2   bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, err := ValidateWif(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateWif() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateWif() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ValidateWif() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("ValidateWif() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func TestToWif(t *testing.T) {
	type args struct {
		prefixHex string
		secretHex string
		compress  bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToWif(tt.args.prefixHex, tt.args.secretHex, tt.args.compress); got != tt.want {
				t.Errorf("ToWif() = %v, want %v", got, tt.want)
			}
		})
	}
}
