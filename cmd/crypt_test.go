package cmd

import (
	"bytes"
	"flag"
	"os"
	"strings"
	"testing"
)

func Test_CryptMain(t *testing.T) {
	const (
		cliName         = "crypt"
		validDecryptOut = "Decrypted message\nTest Message for crypt\n"
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
			name: "valid encrypt",
			args: []string{
				"--address",
				"133r6sCjLq6NbmSLjxuypDuSeocwenu1Qg",
				"--message",
				"Hello World",
				"--signature",
				"Hz0pDxS09fEdjKrSJenLYsz05gakT6eW9GdzKDopnlwUL3bf3mQeJcS+takCIMuavK/NLOorXWXZCqV0KBMlwgU=", "-e",
				"-i",
				"Test Message for crypt",
				"-o",
			},
			want:    0,
			wantOut: "", // Note: Encrypted output changes everytime so we'll just check the prefix for success
		}, {
			name: "valid encrypt",
			args: []string{
				"-i",
				"BFHzllfvzCZbKFXMnTUKitlPlAqiuKXEvs2PopPKx205bZS0GHdvmUaAG2p0R9aBJ3rSiHXrmG4DY7SZS3BKuRyj8Udv2thl/zdAkbuNjs1q98i6FPHLIkAsOaTveAH8cFsFlcwEAZIeA9ExqdpoNhyIU01yS0E=",
				"-k",
				"L165TWkVszAp4yHkFsVRj8udU6w2UxfvVMk8bs9QZZyzNmwWVprK",
				"-d",
				"-o",
			},
			want:    0,
			wantOut: validDecryptOut,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// reset flags else panic
			flag.CommandLine = flag.NewFlagSet(cliName, flag.ExitOnError)
			os.Args = append([]string{cliName}, tt.args...)

			out := &bytes.Buffer{}
			if got := CryptMain(out); got != tt.want {
				t.Errorf("CryptMain() = %v, want %v", got, tt.want)
			}

			gotOut := out.String()
			if !strings.HasPrefix(gotOut, "Encrypted message for 133r6sCj") && gotOut != tt.wantOut {
				t.Errorf("CryptMain() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
