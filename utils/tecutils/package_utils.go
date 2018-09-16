package tecutils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func getPathSeparator() string {
	if runtime.GOOS == "windows" {
		return "\\"
	} else {
		return "/"
	}
}

//Devuelve la ruta local de un package de go
func GetPackageFullPath(name string) (result string, err error) {
	goPath := os.Getenv("GOPATH")
	pkgPath := []string{goPath, "src"}
	for _, p := range strings.Split(name, getPathSeparator()) {
		pkgPath = append(pkgPath, p)
	}
	result = filepath.Join(pkgPath...)
	if ok := DirectoryExists(result); !ok {
		return "", fmt.Errorf("Package not found. Cannot open directory: %s\n", result)
	}

	return result, nil
}
