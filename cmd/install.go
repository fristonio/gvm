package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fristonio/gvm/manager"
	"github.com/fristonio/gvm/network"
	"github.com/fristonio/gvm/utils"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs the version of go mentioned against this flag",
	Long: `Installs the version of go mentioned against this flag
For this it first calls the downloader to download the zip for the version of go
Then install build it to be used`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			utils.Log.Error("No version for go is provided")
			utils.Log.Error(`Use format : gvm install [go version]
    gvm install go1.9
    To list version available for use : gvm list-remote`)
			os.Exit(1)
		}

		releaseName := args[0]
		if utils.GOS_REGEXP.FindString(releaseName) != "" {
			// Once we got go version from the user, check if it already exist in downloads
			// If it does check if it is installed
			// Prompt user to fix it if it is already installed
			// Otherwise download the version source from remote, copy it to Gos directory
			// Build it, create an environment file for it.
			releases, err := network.ParseGoReleases(false)
			if err != nil {
				utils.Log.Errorf("An error occured while parsing available releases : %v", err)
				os.Exit(1)
			}

			var goRelease network.Release
			var flag bool

			for _, release := range releases {
				if release.Name == releaseName {
					goRelease = release
					flag = true
					break
				}
			}

			if !flag {
				utils.Log.Errorf(`Could not find a matching go version source.
	Use gvm list-remote to list all the available versions.`)
				os.Exit(1)
			}

			manageReleaseDownload(goRelease)
			manageCompressedDownload(goRelease)

			// Compile the source of go obtained
			utils.Log.Info("Compiling go from source")
			err = manager.CompileGoRelease(goRelease.Name)
			if err != nil {
				utils.Log.Errorf("Error during compilation : %v", err)
			}
			manager.CreateEnvironmentFile(goRelease.Name)
			os.Exit(0)
		} else {
			utils.Log.Error("Not a valid go version, it should be of the format : goX.X")
			os.Exit(1)
		}
	},
}

func forceNewDownload() bool {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("[*] Already file exist in .gvm/downloads, force download clearing previous files[Y/N] : ")
	scanner.Scan()
	text := scanner.Text()
	if text == "Y" || text == "y" || text == "" {
		return true
	}
	return false
}

func manageReleaseDownload(goRelease network.Release) {
	downloadPath := filepath.Join(utils.GVM_ROOT_DIR, utils.GVM_DOWNLOAD_DIR, filepath.Base(goRelease.DownloadUrl))
	if !utils.CheckIfAlreadyExist(downloadPath) {
		utils.Log.Infof("Beggining to download source for %s", goRelease.Name)
		if err := network.Download(goRelease.DownloadUrl, true, 4, false); err != nil {
			if forceNewDownload() {
				if e := network.Download(goRelease.DownloadUrl, true, 4, true); e != nil {
					utils.Log.Error("An error occured while downloading go from source")
					os.Exit(1)
				}
			} else {
				utils.Log.Error("Error while downloading go version source")
				os.Exit(1)
			}
		}
		utils.Log.Info("Download completed...")
	} else {
		utils.Log.Infof("Found a cached copy for %s", goRelease.Name)
	}
}

func manageCompressedDownload(goRelease network.Release) {
	utils.Log.Info("Unzipping the downloaded source ...")
	source := filepath.Join(
		utils.GVM_ROOT_DIR,
		utils.GVM_DOWNLOAD_DIR,
		filepath.Base(goRelease.DownloadUrl),
	)
	destination := filepath.Join(
		utils.GVM_ROOT_DIR,
		utils.GVM_GOS_DIRNAME,
		goRelease.Name,
	)

	if utils.CheckIfAlreadyExist(source) {
		err := utils.UntarToDestination(source, destination)
		if err != nil {
			utils.Log.Infof("Error while trying to decompress source : %v", err)
			os.Exit(1)
		}
	}
}
