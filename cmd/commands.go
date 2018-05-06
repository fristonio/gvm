package cmd

import (
	"fmt"
	"os"

	"github.com/fristonio/gvm/logger"
	"github.com/spf13/cobra"
)

var longDescriptionGvm = `
A Fast and Flexible version manager for Golang built with
love by fristonio in Go.
Complete source code is available at https://github.com/fristonio/gvm`

var log *logger.Logger = logger.New(os.Stdout)

var rootCmd = &cobra.Command{
	Use:   "gvm",
	Short: "gvm is a fast and reliable version manager for go",
	Long:  longDescriptionGvm,
	Run: func(cmd *cobra.Command, args []string) {
		log.Warn("No arguments are supplied ... ")
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(listRemoteCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(uninstallCmd)
}
