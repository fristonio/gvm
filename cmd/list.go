package cmd

import (
	"io/ioutil"
	"path/filepath"

	"github.com/fristonio/gvm/utils"
	"github.com/spf13/cobra"
)

// Print the version of gvm running and exit gracefully
// Version information are retrived from the version/version.go which is populated at
// build time using ldflags
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List local version of go available for use",
	Long: `List all the releases of golang that are available in the local
gvm environment to use.`,

	Run: func(cmd *cobra.Command, args []string) {
		gosDir := filepath.Join(utils.GVM_ROOT_DIR, utils.GVM_GOS_DIRNAME)
		gos, err := ioutil.ReadDir(gosDir)
		if err != nil {
			log.Fatal("No gos installed, to view a list of versions available use: go list-remote")
		}

		var installedGos = make([]string, 0)
		for _, f := range gos {
			if f.IsDir() && utils.GOS_REGEXP.FindString(f.Name()) != "" {
				installedGos = append(installedGos, f.Name())
			}
		}
		utils.PrintInstalledGos(installedGos)
	},
}
