package main

import (
	"fmt"
	"github.com/nikitatroshenko/wget/utils"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	rawurl := os.Args[1]
	parsed, err := url.Parse(rawurl)
	check(err)
	log.Printf("scheme: %s\n", parsed.Scheme)

	resp, err := http.Get(rawurl)
	check(err)

	body, err := ioutil.ReadAll(resp.Body)
	check(err)
	err = resp.Body.Close()
	check(err)
	filename := utils.UrlFileName(parsed)
	log.Printf("filename: '%s'\n", filename)
	file, _ := os.Create(filename)
	defer file.Close()
	_, err = fmt.Fprintf(file, "%s", body)

	check(err)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
