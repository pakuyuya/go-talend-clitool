package fileutils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func FindMatchPathes(path string) []string {
	if os.PathSeparator != '/' {
		path = strings.Replace(path, "\\", "/", -1)
	}

	basepath := "."
	if filepath.IsAbs(path) || path[0:1] == "/" {
		var pos int
		basepath, pos = _fetchToken(path)
		path = path[pos:]
	}

	return _findMatchPaths(basepath, path)
}

func _findMatchPaths(basepath, path string) []string {
	rets := []string{}

	if path == "" {
		return make([]string, 0, 0)
	}

	t, i := _fetchToken(path)
	nextpath := path[i:]

	if t == "**" {
		return _findMatchPathsRecursive(basepath, nextpath)
	}

	regexptn := strings.Replace(regexp.QuoteMeta(t), "\\*", ".*", -1)
	r := regexp.MustCompile(regexptn)

	files, _ := ioutil.ReadDir(basepath)
	for _, file := range files {
		childfile := file.Name()
		if r.MatchString(childfile) {
			if nextpath == "" {
				rets = append(rets, basepath+file.Name())
			} else if file.IsDir() {
				rets = append(rets, _findMatchPaths(basepath+file.Name(), nextpath)...)
			}

		}
	}
	return rets
}

func _findMatchPathsRecursive(basepath, path string) []string {
	if path == "" {
		return []string{}
	}
	if path == "**" {
		_, i := _fetchToken(path)
		return _findMatchPathsRecursive(basepath, path[i:])
	}

	rets := []string{}

	t, i := _fetchToken(path)
	regexptn := strings.Replace(regexp.QuoteMeta(t), "\\*", ".*", -1)
	r := regexp.MustCompile(regexptn)

	files, _ := ioutil.ReadDir(basepath)
	for _, file := range files {
		childfile := file.Name()
		if r.MatchString(childfile) {
			if file.IsDir() {
				rets = append(rets, _findMatchPaths(basepath+"/"+file.Name(), path[i:])...)
			} else {
				rets = append(rets, basepath+"/"+file.Name())
			}
		} else {
			if file.IsDir() {
				rets = append(rets, _findMatchPathsRecursive(basepath+"/"+childfile, path)...)
			}
		}
	}

	return rets
}

func _fetchToken(path string) (string, int) {
	pos := _min(strings.Index(path, "/")+1, len(path))
	if pos <= 0 {
		pos = len(path)
	}
	return path[0:pos], pos
}

func _min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
