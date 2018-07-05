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
		path = path[pos+1:]
	}

	return _findMatchPaths(basepath, path)
}

func _findMatchPaths(basepath, path string) []string {
	rets := []string{}

	if path == "" {
		return make([]string, 0, 0)
	}

	t, i := _fetchToken(path)
	path = path[i+1:]

	if t == "**" {
		return _findMatchPathsRecursive(basepath, path)
	}

	regexptn := strings.Replace(regexp.QuoteMeta(t), "\\*", ".*", -1)
	r := regexp.MustCompile(regexptn)

	files, _ := ioutil.ReadDir(basepath)
	for _, file := range files {
		childfile := file.Name()
		if r.MatchString(childfile) {
			if file.IsDir() {
				rets = append(rets, _findMatchPaths(basepath+"/"+file.Name(), path)...)
			} else {
				rets = append(rets, basepath+"/"+file.Name())
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
		path, _ = _fetchToken(path)
		return _findMatchPathsRecursive(basepath, path)
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
				rets = append(rets, _findMatchPaths(basepath+"/"+file.Name(), path[i+1:])...)
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
	pos := _min(strings.Index(path, "/"), len(path))
	return path[0:pos], pos
}

func _min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
