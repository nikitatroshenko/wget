package http

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var method = flag.String("method", "", "HTTP method to send to server")
var bodyData = flag.String("body-data", "", "String that needs to be send along with Method specified using --method. Cannot be used with --body-file")
var bodyFile = flag.String("body-file", "", "Name of file from which to take request body along with Method specified using --method. Cannot be used with --body-data")
var postData = flag.String("post-data", "", "String that needs to be send along with POST Method. Since POST method is implied here, this flag cannot be combined with --method. Cannot be used with --post-file")
var postFile = flag.String("post-file", "", "Name of file from which to take request body  along with POST Method. Since POST method is implied here, this flag cannot be combined with --method. Cannot be used with --post-data")

// WgetHttpResource performs HTTP request to the specified URL and saves response to the provided writer
// HTTP method is passed through command line arguments
func WgetHttpResource(rawurl string) (io.Reader, error) {
	req, err := newRequest(rawurl)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return nil, err
	}
	return &buf, nil
}

type flagValuesMisuseError struct {
	what string
}

func (e *flagValuesMisuseError) Error() string {
	return fmt.Sprintf("%s", e.what)
}

func newRequest(rawurl string) (*http.Request, error) {
	if *method != "" {
		if *postData != "" || *postFile != "" {
			return nil, &flagValuesMisuseError{"You cannot use --post-file or --post-data along with --method. --method expects data through --body-data or --body-file options"}
		}
		body, err := getBody()
		if err != nil {
			return nil, err
		}
		return http.NewRequest(*method, rawurl, body)
	}
	if *bodyData != "" || *bodyFile != "" {
		return nil, &flagValuesMisuseError{"You must specify a method through --method=HTTPMethod to use with --body-data or --body-file"}
	}
	if *postData != "" || *postFile != "" {
		body, err := getPostData()
		if err != nil {
			return nil, err
		}
		return http.NewRequest("POST", rawurl, body)
	}
	return http.NewRequest("GET", rawurl, nil)
}

func getBody() (io.Reader, error) {
	if *bodyData != "" && *bodyFile != "" {
		return nil, &flagValuesMisuseError{"You cannot specify both --body-file and --body-data at the same time"}
	}
	if *bodyData != "" {
		return strings.NewReader(*bodyData), nil
	}
	if *bodyFile != "" {
		return os.Open(*bodyFile)
	}
	return nil, nil
}

func getPostData() (io.Reader, error) {
	if *postData != "" && *postFile != "" {
		return nil, &flagValuesMisuseError{"You cannot specify both --post-file and --post-data at the same time"}
	}
	if *postData != "" {
		return strings.NewReader(*postData), nil
	}
	if *postFile != "" {
		return os.Open(*postFile)
	}
	return nil, nil
}
