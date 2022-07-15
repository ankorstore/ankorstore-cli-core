package util

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func IsBrewInstallation() bool {
	ankorPath, _ := exec.LookPath(AppName)
	if strings.Contains(ankorPath, Homebrew) {
		return true
	}
	return false
}

type PluginName struct {
	Name string
	Type string
}

type PluginInfo struct {
	Name         *PluginName
	Dependencies []*PluginName
	Fields       map[string]string
}

func PathByName(name string) string {
	pluginDir := []string{
		"", // path to plugin dir
		fmt.Sprintf("ankor-%s_%s_%s", name, runtime.GOOS, runtime.GOARCH),
	}
	path := pluginDir
	return filepath.Join(path...)
}
