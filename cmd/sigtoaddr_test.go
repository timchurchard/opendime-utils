package cmd

import (
	"bytes"
	"flag"
	"os"
	"testing"
)

func Test_sanityTest(t *testing.T) {
	if sanityTests() != 0 {
		t.Fail()
	}
}

func Test_SigtoaddrMain(t *testing.T) {
	const (
		cliName                 = "sigtoaddr"
		bitcoinValidExpectedout = "Addresses for Opendime:\t1Nu1QpfegiGmqHS6YZxkaiGpnqAUXvZz2f\n- Bitcoin P2PKH\t\t\t 1GLfgL9yKVTRRG1D4fdKkEuEQqAE7ob1eB \n- Bitcoin P2PKH (Compressed)\t 1Nu1QpfegiGmqHS6YZxkaiGpnqAUXvZz2f \n- Bitcoin P2WPKH\t\t bc1q7qcf63rtp20dsalcwmceucxs0kwn75l95nsxjf \n- Ethereum\t\t\t 0x5D0a9F69035Be4275204f9eBbd5cC049e42429c6 \n- Litecoin P2PKH\t\t LaZcwYToQ9hUg4hNEocd2Fxzd3XWEMFnQ5 \n- Litecoin P2PKH (Compressed)\t Lh7xg2yUmNWq668Fihx3rjLb13XkbHuBMQ \n- Litecoin P2WPKH\t\t ltc1q7qcf63rtp20dsalcwmceucxs0kwn75l9s02z2e \n- Dogecoin P2PKH\t\t DLUmDb6ccuMhxGBooFctJ14qHxtXVWYf4P \n"
	)

	// We manipulate the Args to set them up for the testcases, after this test we restore the initial args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	tests := []struct {
		name    string
		args    []string
		want    int
		wantOut string
	}{
		{
			name: "valid bitcoin",
			args: []string{
				"--address",
				"1Nu1QpfegiGmqHS6YZxkaiGpnqAUXvZz2f",
				"--signature",
				"HwPlEOxTxs62ruMHZvamv0wmUlbbaY/2ZSqw9Hpdw+FWfgXuSxQ9x55ceSiFyvnlpiZjt+KIhSYnhGnCv8iDe5o=",
				"--message",
				"Hello World!",
			},
			want:    0,
			wantOut: bitcoinValidExpectedout,
		}, {
			name: "invalid bitcoin",
			args: []string{
				"--address",
				"1Nu1QpfegiGmqHS6YZxkaiGpnqAUXvZz2f",
				"--signature",
				"HwPlEOxTxs62ruMHZvamv0wmUlbbaY/2ZSqw9Hpdw+FWfgXuSxQ9x55ceSiFyvnlpiZjt+KIhSYnhGnCv8iDe5o=",
				"--message",
				"invalid invalid invalid",
			},
			want:    1,
			wantOut: "Unable to verify signature: Invalid signature address not match",
		},
	}
	for _, tt := range tests {
		// reset flags else panic
		flag.CommandLine = flag.NewFlagSet(cliName, flag.ExitOnError)
		os.Args = append([]string{cliName}, tt.args...)

		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			if got := SigtoaddrMain(out); got != tt.want {
				t.Errorf("SigtoaddrMain() = %v, want %v", got, tt.want)
			}
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("SigtoaddrMain() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
