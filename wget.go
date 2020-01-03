package main

import (
	"flag"
	"fmt"
	"github.com/nikitatroshenko/wget/http"
	"github.com/nikitatroshenko/wget/utils"
	"io"
	"log"
	"net/url"
	"os"
)

func main() {
	flag.Parse()

	tailArgs := flag.Args()
	if len(tailArgs) == 0 {
		log.Fatal("wget: missing URL")
	}
	rawurl := tailArgs[0]
	parsed, err := parseURL(rawurl)
	check(err)

	fileName := utils.UrlFileName(parsed)
	log.Printf("fileName: '%s'\n", fileName)

	response, err := http.WgetHttpResource(rawurl)
	check(err)
	file, err := os.Create(fileName)
	check(err)
	defer file.Close()
	n, err := io.Copy(file, response)
	log.Printf("'%s' saved [%d]", fileName, n)
}

type unsupportedSchemeError struct {
	url    string
	scheme string
}

func (e *unsupportedSchemeError) Error() string {
	return fmt.Sprintf("%s: Unsupported scheme '%s'", e.url, e.scheme)
}

// parseURL uses net/url.Parse function but with additional validation of supported scheme
func parseURL(rawurl string) (*url.URL, error) {
	parsed, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	if !isSchemeSupported(parsed.Scheme) {
		return nil, &unsupportedSchemeError{rawurl, parsed.Scheme}
	}
	return parsed, nil
}

func isSchemeSupported(scheme string) bool {
	return scheme == "http"
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
