package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/timchurchard/opendime-utils/cmd"
)

const cliName = "opendime-utils"

func main() {
	if len(os.Args) < 2 {
		usageRoot()
	}

	// Save the command and reset the flags
	command := os.Args[1]
	flag.CommandLine = flag.NewFlagSet(cliName, flag.ExitOnError)
	os.Args = append([]string{cliName}, os.Args[2:]...)

	switch command {
	case "sigtoaddr":
		os.Exit(cmd.SigtoaddrMain(os.Stdout))
	case "keyconv":
		var key string

		fmt.Printf("Private Key WIF: ")
		fmt.Scanln(&key)
		key = strings.TrimSpace(key)

		os.Exit(cmd.KeyconvMain(os.Stdout, key))
	case "crypt":
		os.Exit(cmd.CryptMain(os.Stdout))
	}

	usageRoot()
}

func usageRoot() {
	fmt.Printf("usage: opendime-utils command(sigtoaddr|keyconv|crypt) options\n")
	os.Exit(1)
}
