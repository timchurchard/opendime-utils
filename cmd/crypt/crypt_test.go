package main

import (
	"bytes"
	"flag"
	"os"
	"strings"
	"testing"
)

func Test_realMain(t *testing.T) {
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
				"BEWs+6LFlRxUQQ4X9kuPp03L3C+qttMjliLWVExapsDQdjaFfV/7sPxHbhVDPeZ2upYx99TzK0TufWEupSUAXC7s69dbdqUTiTYZOkfKRahrpJNmTffwbmgIO+lI8qNk/SBXVR/CNu+toq/H+5KqJ6njeqaNZX8=",
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
			if got := realMain(out); got != tt.want {
				t.Errorf("realMain() = %v, want %v", got, tt.want)
			}

			gotOut := out.String()
			if !strings.HasPrefix(gotOut, "Encrypted message for 133r6sCj") && gotOut != tt.wantOut {
				t.Errorf("realMain() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
