package pkg

import (
	"testing"
)

// TestCheckBalance tests against a real api TODO FIXME
func TestCheckBalance(t *testing.T) {
	type args struct {
		address string
		fiat    string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		want1   float64
		wantErr bool
	}{
		{"valid bitcoin", args{address: "1FHxL2JskCy6g98wEMxpaNNkxohjq3hUKk", fiat: "usd"}, 0.01, 100, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := CheckBalance(tt.args.address, tt.args.fiat)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Note: Using <= here because calling real API
			if got <= tt.want {
				t.Errorf("CheckBalance() got = %v, want %v", got, tt.want)
			}
			if got1 <= tt.want1 {
				t.Errorf("CheckBalance() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
