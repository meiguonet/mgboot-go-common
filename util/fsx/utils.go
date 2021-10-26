package fsx

import (
	"path/filepath"
	"runtime"
	"strings"
)

func IsWin() bool {
	s1 := strings.ToLower(runtime.GOOS)
	return !strings.Contains(s1, "darwin") && strings.Contains(s1, "win")
}

func GetExtension(fpath string) string {
	if !strings.Contains(fpath, ".") {
		return ""
	}

	return strings.ToLower(filepath.Ext(fpath))
}

func GetRealpath(basedir, fpath string) string {
	if basedir == "" || !strings.HasPrefix(fpath, "classpath:") {
		return fpath
	}

	fpath = strings.TrimPrefix(fpath, "classpath:")
	fpath = strings.ReplaceAll(fpath, "\\", "/")
	fpath = strings.Trim(fpath, "/")
	basedir = strings.ReplaceAll(basedir, "\\", "/")
	basedir = strings.TrimSuffix(basedir, "/")
	return basedir + "/" + fpath
}
