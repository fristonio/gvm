package utils

import (
	"os"
	"path/filepath"
)

const (
	GVM_DOWNLOAD_DIR    string = "downloads"
	GVM_GOS_DIRNAME     string = "gos"
	GVM_ENV_DIRNAME     string = "environment"
	GVM_PKGSET_NAME     string = "global"
	GVM_PKGSET_DIRNAME  string = "pkgsets"
	GVM_OVERLAY_DIRNAME string = "overlay"
)

var (
	GVM_ROOT_DIR string = filepath.Join(os.Getenv("HOME"), ".gvm")
)

var ENV_FILE string = `#!/bin/bash
# Auto generated shell script to enable an environment for gos

export GVM_ROOT="%s"
export GVM_GO_VERSION="%s"
export GVM_PACKAGESET_NAME="%s"
export GOROOT="%s"
export GOPATH="%s"
export GVM_OVERLAY_PREFIX="%s"
export PATH="%s"
export LD_LIBRARY_PATH="%s"
export DYLD_LIBRARY_PATH="%s"
export PKG_CONFIG_PATH="%s"

`
