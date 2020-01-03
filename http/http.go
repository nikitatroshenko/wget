package http

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var method = flag.String("method", "GET", "HTTP method to send to server")
var bodyData = flag.String("body-data", "", "String that needs to be send along with Method specified using --method. Cannot be used with --body-file")
var bodyFile = flag.String("body-file", "", "Name of file from which to take request body. Cannot be used with --body-data")

// WgetHttpResource performs HTTP request to the specified URL and saves response to the provided writer
// HTTP method is passed through command line arguments
func WgetHttpResource(rawurl string, writer io.Writer) (written int64, err error) {
	requestBody, err := getBody()
	if err != nil {
		return 0, err
	}
	req, err := http.NewRequest(*method, rawurl, requestBody)
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

type flagValuesConflictError struct {
	what string
}

func (e *flagValuesConflictError) Error() string {
	return fmt.Sprintf("%s", e.what)
}

func getBody() (io.Reader, error) {
	if *bodyData != "" && *bodyFile != "" {
		return nil, &flagValuesConflictError{"You cannot use specify --body-file and --body-data at the same time"}
	}
	if *bodyData != "" {
		return strings.NewReader(*bodyData), nil
	}
	if *bodyFile != "" {
		return os.Open(*bodyFile)
	}
	return nil, nil
}
