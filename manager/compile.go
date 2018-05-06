package manager

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fristonio/gvm/utils"
)

func CompileGoRelease(releaseName string) error {
	err := CreateCompilationEnv(releaseName)
	if err != nil {
		return err
	}
	goSrcDir := filepath.Join(utils.GVM_ROOT_DIR, utils.GVM_GOS_DIRNAME, releaseName, "src")
	os.Chdir(goSrcDir)

	cmd := exec.Command("./make.bash")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("Error while running compilation : %v\n", err)
	}
	return nil
}
