package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fristonio/gvm/logger"
	"github.com/fristonio/gvm/version"
)

var log *logger.Logger = logger.New(os.Stdout)
var versionFlag = flag.Bool("version", false, "Prints gvm version and exit")

func main() {
	flag.Parse()
	args := flag.Args()

	// Print the version of gvm running and exit gracefully
	// Version information are retrived from the version/version.go which is populated at
	// build time using ldflags
	if *versionFlag {
		fmt.Printf("gvm version : %s(%s)\n", version.Info["version"], version.Info["branch"])
		os.Exit(0)
	}

	if len(args) == 0 {
		log.Error("No arguments are supplied")
		flag.Usage()
	}
}
