package cmd

import (
	"github.com/fristonio/gvm/network"
	"github.com/spf13/cobra"
)

// Print the version of gvm running and exit gracefully
// Version information are retrived from the version/version.go which is populated at
// build time using ldflags
var listRemoteCmd = &cobra.Command{
	Use:   "list-remote",
	Short: "List remote version of go available",
	Long:  `List all the releases of golang that are available`,

	Run: func(cmd *cobra.Command, args []string) {
		_, err := network.ParseGoReleases(true)
		if err != nil {
			log.Errorf("An error occured while parsing available releases : %v", err)
		}
	},
}
