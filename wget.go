package main

import (
	"flag"
	"fmt"
	"github.com/nikitatroshenko/wget/utils"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	flag.Parse()

	tailArgs := flag.Args()
	rawurl := tailArgs[0]
	parsed, err := url.Parse(rawurl)
	check(err)
	log.Printf("scheme: %s\n", parsed.Scheme)
	filename := utils.UrlFileName(parsed)
	log.Printf("filename: '%s'\n", filename)

	file, _ := os.Create(filename)
	defer file.Close()
	wgetHttpResource(rawurl, file)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func wgetHttpResource(rawurl string, writer io.Writer) {
	resp, err := http.Get(rawurl)
	check(err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	check(err)
	_, err = fmt.Fprintf(writer, "%s", body)
	check(err)
}
