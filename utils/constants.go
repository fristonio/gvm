package utils

import "path/filepath"

const (
	GVM_DOWNLOAD_DIR string = "downloads"
	GVM_GOS_DIRNAME  string = "gos"
	GVM_ENV_DIRNAME  string = "environment"
)

var (
	GVM_ROOT_DIR string = filepath.Join(os.Getenv("HOME"), ".gvm")
)
