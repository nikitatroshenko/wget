package http

import (
	"flag"
	"io"
	"net/http"
	"strings"
)

var method = flag.String("method", "GET", "HTTP method to send to server")
var bodyData = flag.String("body-data", "", "String that needs to be send along with Method specified using --method")

// WgetHttpResource performs HTTP request to the specified URL and saves response to the provided writer
// HTTP method is passed through command line arguments
func WgetHttpResource(rawurl string, writer io.Writer) (written int64, err error) {
	req, err := http.NewRequest(*method, rawurl, strings.NewReader(*bodyData))
	if err != nil {
		return 0, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	n, err := io.Copy(writer, resp.Body)
	return n, err
}
