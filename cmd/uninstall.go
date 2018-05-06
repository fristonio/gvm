package cmd

import (
	"os"
	"path/filepath"

	"github.com/fristonio/gvm/utils"
	"github.com/spf13/cobra"
)

// Print the version of gvm running and exit gracefully
// Version information are retrived from the version/version.go which is populated at
// build time using ldflags
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall the specified version of go",
	Long:  `Uninstall, remove the env file and delete the source directory for the version specified as the argument to this command`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			utils.Log.Error("No version for go is provided")
			utils.Log.Error(`Use format : gvm use [go version]
    gvm use go1.9
    To list version available for use : gvm list`)
			os.Exit(1)
		}

		releaseName := args[0]
		if utils.GOS_REGEXP.FindString(releaseName) != "" {
			envFile := filepath.Join(utils.GVM_ROOT_DIR, utils.GVM_ENV_DIRNAME, releaseName)
			if _, err := os.Stat(envFile); !os.IsNotExist(err) {
				os.Remove(envFile)
			}

			goSrcDir := filepath.Join(utils.GVM_ROOT_DIR, utils.GVM_GOS_DIRNAME, releaseName)
			if _, err := os.Stat(goSrcDir); !os.IsNotExist(err) {
				os.RemoveAll(goSrcDir)
			}

			goPkgsetDir := filepath.Join(utils.GVM_ROOT_DIR, utils.GVM_PKGSET_DIRNAME, releaseName)
			if _, err := os.Stat(goPkgsetDir); !os.IsNotExist(err) {
				os.RemoveAll(goPkgsetDir)
			}
		} else {
			utils.Log.Error("Not a valid go version, it should be of the format : goX.X")
			os.Exit(1)
		}
	},
}
