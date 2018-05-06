package manager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fristonio/gvm/utils"
)

// Create global package directories for GVM
// It includes Creating the following directory structure for a given go version
/*
[fristonio] $ /tmp/.gvm
❮❮ tree
.
├── downloads
├── environments
├── gos
	└── go1.9
└── pkgsets
    └── go1.9
        └── global
            └── overlay
                ├── bin
                └── lib
                    └── pkgconfig

10 directories, 0 files
*/
func CreateGlobalPackageSets(goVersion string) error {
	gvmPkgSet := filepath.Join(utils.GVM_ROOT_DIR, utils.GVM_PKGSET_DIRNAME, goVersion)
	err := utils.CreateDirStrucutre(gvmPkgSet)
	if err != nil {
		return fmt.Errorf("Error while creating GVM packageset")
	}

	gvmOverlayRoot := filepath.Join(gvmPkgSet, utils.GVM_PKGSET_NAME, utils.GVM_OVERLAY_DIRNAME)

	gvmOverlayPkgConfig := filepath.Join(gvmOverlayRoot, "lib", "pkgconfig")
	err = utils.CreateDirStrucutre(gvmOverlayPkgConfig)
	if err != nil {
		return fmt.Errorf("Error while creating GVM packageset overlay")
	}

	gvmOverlayBin := filepath.Join(gvmOverlayRoot, "bin")
	err = utils.CreateDirStrucutre(gvmOverlayBin)
	if err != nil {
		return fmt.Errorf("Error while creating GVM packageset overlay/bin")
	}
	return nil
}

// Create environment file for activating a go version in golang
// Each go version in associated  with an environment shell script
// Which creates the required environment for that version of go
// version specifies the version of the golang we are creating the env for
func CreateEnvironmentFile(goVersion string) error {
	CreateGlobalPackageSets(goVersion)
	if utils.GOS_REGEXP.FindString(goVersion) == "" {
		errStr := fmt.Sprintf("Not a valid go name %s to create environment", goVersion)
		utils.Log.Warn(errStr)
		return fmt.Errorf(errStr)
	}
	var environmentDir string = filepath.Join(utils.GVM_ROOT_DIR, utils.GVM_ENV_DIRNAME)
	err := utils.CreateDirStrucutre(environmentDir)
	if err != nil {
		utils.Log.Warnf("An error occured while creating enviroment directory : %s", environmentDir)
		return err
	}
	var environmentFile string = filepath.Join(environmentDir, goVersion)
	_, err = os.Stat(environmentFile)
	if err == nil {
		// Environment already exist, so just to be on the safe side create the environment again
		// With the latest information we have about it.
		os.Remove(environmentFile)
	}

	err = nil
	file, err := os.OpenFile(environmentFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0775)
	if err != nil {
		errStr := fmt.Sprintf("An error occured while opening file : %s", environmentFile)
		return fmt.Errorf(errStr)
	}
	defer file.Close()

	// Get the variables ready for environment file
	gvmGosRoot := filepath.Join(utils.GVM_ROOT_DIR, utils.GVM_GOS_DIRNAME, goVersion)
	gvmGoPath := filepath.Join(utils.GVM_ROOT_DIR, utils.GVM_PKGSET_DIRNAME, goVersion, utils.GVM_PKGSET_NAME)
	gvmGoOverlayPath := filepath.Join(gvmGoPath, utils.GVM_OVERLAY_DIRNAME)

	gvmGosRootBin := filepath.Join(gvmGosRoot, "bin")
	gvmGoBinPath := filepath.Join(gvmGoPath, "bin")
	gvmGoOverlayBinPath := filepath.Join(gvmGoOverlayPath, "bin")
	gvmRootBinPath := filepath.Join(utils.GVM_ROOT_DIR, "bin")
	newENVPath := gvmGosRootBin + ":" + gvmGoBinPath + ":" + gvmGoOverlayBinPath + ":" + gvmRootBinPath + `:$PATH`

	gvmOverlayLibPath := filepath.Join(gvmGoOverlayPath, "lib")
	newLdLibPath := gvmOverlayLibPath + ":$LD_LIBRARY_PATH"
	newDyldLibPath := gvmOverlayLibPath + ":$DYLD_LIBRARY_PATH"

	gvmOverlayPkgConfig := filepath.Join(gvmOverlayLibPath, "pkgconfig")
	newPkgConfigPath := gvmOverlayPkgConfig + ":$PKG_CONFIG_PATH"

	goEnv := fmt.Sprintf(utils.ENV_FILE,
		utils.GVM_ROOT_DIR,
		goVersion,
		utils.GVM_PKGSET_NAME,
		gvmGosRoot,
		gvmGoPath,
		gvmGoOverlayPath,
		newENVPath,
		newLdLibPath,
		newDyldLibPath,
		newPkgConfigPath,
	)

	_, e := file.WriteString(goEnv)
	if e != nil {
		return fmt.Errorf("An error occured while writing the environment configuration file")
	}
	return nil
}

// Create environment for compilation of go from source
// Unsets previously set env variable and set to new ones.
// Take a look at manager/new_installation.md to get an insight for the procedure
func CreateCompilationEnv(goVersion string) error {
	var pathEnvVar string = os.Getenv("PATH")
	var goVerDir string = filepath.Join(utils.GVM_ROOT_DIR, utils.GVM_GOS_DIRNAME, goVersion)
	var gobinEnvPath string = filepath.Join(goVerDir, "bin")

	err := utils.CheckIfDirExist(goVerDir)
	if err != nil {
		errStr := fmt.Sprintf("No sources with goVersion %s found.", goVersion)
		return fmt.Errorf(errStr)
	}

	var baseGoRoot string = os.Getenv("GOROOT")
	if baseGoRoot == "" {
		return fmt.Errorf("No GOROOT environment variable")
	}
	os.Setenv("GOROOT_BOOTSTRAP", baseGoRoot)
	os.Unsetenv("GOARCH")
	os.Unsetenv("GOOS")
	os.Unsetenv("GOPATH")
	os.Unsetenv("GOBIN")
	os.Unsetenv("GOROOT")

	os.Setenv("GOBIN", gobinEnvPath)
	os.Setenv("PATH", gobinEnvPath+pathEnvVar)
	os.Setenv("GOROOT", goVerDir)
	return nil
}
