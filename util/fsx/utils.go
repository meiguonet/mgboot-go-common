package fsx

import (
	"github.com/meiguonet/mgboot-go-common/AppConf"
	"os"
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

	ext := strings.ToLower(filepath.Ext(fpath))

	if strings.HasPrefix(ext, ".") {
		return ext[1:]
	}

	return ext
}

func GetRealpath(fpath string, baseDir ...string) string {
	fpath = strings.ReplaceAll(fpath, "\\", "/")
	fpath = strings.TrimRight(fpath, "/")
	var dir string

	if len(baseDir) > 0 {
		dir = baseDir[0]
	}

	if dir == "" {
		dir = AppConf.GetDataDir()
	}

	dir = strings.ReplaceAll(dir, "\\", "/")
	dir = strings.TrimRight(dir, "/")
	var relative bool
	var relativeWd bool

	if strings.HasPrefix(fpath, "classpath:") {
		relative = true
		fpath = strings.TrimPrefix(fpath, "classpath:")
		fpath = strings.TrimLeft(fpath, "/")
	} else if strings.HasPrefix(fpath, "{classpath}") {
		relative = true
		fpath = strings.TrimPrefix(fpath, "{classpath}")
		fpath = strings.TrimLeft(fpath, "/")
	} else if strings.HasPrefix(fpath, "datadir:") {
		relative = true
		fpath = strings.TrimPrefix(fpath, "datadir:")
		fpath = strings.TrimLeft(fpath, "/")
	} else if strings.HasPrefix(fpath, "{datadir}") {
		relative = true
		fpath = strings.TrimPrefix(fpath, "{datadir}")
		fpath = strings.TrimLeft(fpath, "/")
	} else if strings.HasPrefix(fpath, "ProjectRoot:") {
		relative = true
		fpath = strings.TrimPrefix(fpath, "ProjectRoot:")
		fpath = strings.TrimLeft(fpath, "/")
	} else if strings.HasPrefix(fpath, "project_root:") {
		relative = true
		fpath = strings.TrimPrefix(fpath, "project_root:")
		fpath = strings.TrimLeft(fpath, "/")
	} else if strings.HasPrefix(fpath, "{ProjectRoot}") {
		relative = true
		fpath = strings.TrimPrefix(fpath, "{ProjectRoot}")
		fpath = strings.TrimLeft(fpath, "/")
	} else if strings.HasPrefix(fpath, "{project_root}") {
		relative = true
		fpath = strings.TrimPrefix(fpath, "{project_root}")
		fpath = strings.TrimLeft(fpath, "/")
	} else if strings.HasPrefix(fpath, "AppRoot:") {
		relative = true
		fpath = strings.TrimPrefix(fpath, "AppRoot:")
		fpath = strings.TrimLeft(fpath, "/")
	} else if strings.HasPrefix(fpath, "app_root:") {
		relative = true
		fpath = strings.TrimPrefix(fpath, "app_root:")
		fpath = strings.TrimLeft(fpath, "/")
	} else if strings.HasPrefix(fpath, "{AppRoot}") {
		relative = true
		fpath = strings.TrimPrefix(fpath, "{AppRoot}")
		fpath = strings.TrimLeft(fpath, "/")
	} else if strings.HasPrefix(fpath, "{app_root}") {
		relative = true
		fpath = strings.TrimPrefix(fpath, "{app_root}")
		fpath = strings.TrimLeft(fpath, "/")
	} else if strings.HasPrefix(fpath, "wd:") {
		relativeWd = true
		fpath = strings.TrimPrefix(fpath, "wd:")
		fpath = strings.TrimLeft(fpath, "/")
	} else if strings.HasPrefix(fpath, "{wd}") {
		relativeWd = true
		fpath = strings.TrimPrefix(fpath, "{wd}")
		fpath = strings.TrimLeft(fpath, "/")
	} else if strings.HasPrefix(fpath, "cwd:") {
		relativeWd = true
		fpath = strings.TrimPrefix(fpath, "cwd:")
		fpath = strings.TrimLeft(fpath, "/")
	} else if strings.HasPrefix(fpath, "{cwd}") {
		relativeWd = true
		fpath = strings.TrimPrefix(fpath, "{cwd}")
		fpath = strings.TrimLeft(fpath, "/")
	}

	if relativeWd {
		wd, _ := os.Getwd()

		if wd == "" {
			return fpath
		}

		wd = strings.ReplaceAll(wd, "\\", "/")
		wd = strings.TrimRight(wd, "/")
		return wd + "/" + fpath
	}

	if relative {
		if dir == "" {
			return fpath
		}

		dir = strings.ReplaceAll(dir, "\\", "/")
		dir = strings.TrimRight(dir, "/")
		return dir + "/" + fpath
	}

	return fpath
}
