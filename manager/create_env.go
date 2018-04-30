package manager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fristonio/gvm/utils"
)

// Create environment file for activating a go version in golang
// Each go version in associated  with an environment shell script
// Which creates the required environment for that version of go
// version specifies the version of the golang we are creating the env for
func CreateEnvironmentFile(goVersion string) error {
	if utils.GOS_REGEXP.FindString(goVersion) == "" {
		errStr := fmt.Sprintf("Not a valid go name %s to create environment", goVersion)
		utils.Log.Warn(errStr)
		return error.New(errStr)
	}
	var environmentDir string = filepath.Join(utils.GVM_ROOT_DIR, utils.GVM_ENV_DIRNAME)
	err := utils.CreateDirStrucutre
	if err != nil {
		utils.Log.Warnf("An error occured while creating enviroment directory : %s", environmentDir)
		return err
	}
	var environmentFile string = filepath.Join(environmentDir, goVersion)
	err := os.Stat(environmentFile)
	if err == nil {
		// Environment already exist, so just to be on the safe side create the environment again
		// With the latest information we have about it.
	}
	return nil
}

// Create environment for compilation of go from source
// Unsets previously set env variable and set to new ones.
// Take a look at manager/new_installation.md to get an insight for the procedure
func CreateCompilationEnv(goVersion string) error {
	var gobinEnvPath string = filePath.Join(goVerDir, "bin")
	var pathEnvVar string = os.Getenv("PATH")
	var goVerDir string = filepath.Join(utils.GVM_ROOT_DIR, utils.GVM_GOS_DIRNAME, goVersion)

	err := utils.CheckIfDirExist(goVerDir)
	if err != nil {
		errStr := fmt.Sprintf("No sources with goVersion %s found.", goVersion)
		return error.New(errStr)
	}

	var baseGoRoot string = os.Getenv("GOROOT")
	if baseGoRoot == "" {
		return error.New("No GOROOT environment variable")
	}
	err := os.Setenv("GOROOT_BOOTSTRAP", baseGoRoot)
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
