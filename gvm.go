package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fristonio/gvm/version"
)

var versionFlag = flag.Bool("version", false, "Prints gvm version and exit")

func init() {
	// Set Defalut logging verbosity to 2
	flag.Set("v", "2")
}

func main() {
	flag.Parse()

	// Print the version of gvm running and exit gracefully
	// Version information are retrived from the version/version.go which is populated at
	// build time using ldflags
	if *versionFlag {
		fmt.Printf("gvm version : %s(%s)\n", version.Info["version"], version.Info["branch"])
		os.Exit(0)
	}
}
