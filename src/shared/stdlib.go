package shared

import (
	_ "embed"
	"os"
	"path/filepath"
)

func getRootDirectory() string {
	dir, _ := os.Executable()
	return dir
}

var TekoRootDirectory string = getRootDirectory()

var StdLibLocation string = filepath.Join(TekoRootDirectory, "lib/stdlib.to")

//go:embed stdlib.to
var StdLibContents string
