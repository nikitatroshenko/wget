package utils

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

const (
	defaultFileName = "index.html"
)

// UrlFileName  takes parsed url.URL and infers local file name from
// its path section. If path is empty then 'index.html' is used
func UrlFileName(urlParsed *url.URL) string {
	urlPath := urlParsed.RawPath

	if urlPath == "" {
		urlPath = urlParsed.Path
	}
	fileName := strings.TrimRight(urlPath, "/")
	if i := strings.LastIndex(fileName, "/"); i >= 0 {
		fileName = fileName[i+1:]
	}
	if fileName == "" {
		fileName = defaultFileName
	}

	return UniqueName(fileName)
}

// UniqueName tries to find a unique name based on prefix provided.
func UniqueName(prefix string) string {
	for count := 0; ; count++ {
		filename := tryUniqueName(prefix, count)
		if filename != "" {
			return filename
		}
	}
}

func tryUniqueName(prefix string, count int) string {
	filename := prefix
	if count > 0 {
		filename = fmt.Sprintf("%s.%d", prefix, count)
	}
	if isFileExist(filename) {
		return ""
	}
	return filename
}

func isFileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
