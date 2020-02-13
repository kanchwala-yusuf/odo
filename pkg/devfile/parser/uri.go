package parser

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/golang/glog"
	"github.com/pkg/errors"
)

const (
	httpsScheme = "https"
	httpScheme  = "http"
	fileScheme  = "file"
)

// fetchFileHttp fetches content from the given http(s) url and stores it
// locally in a temp file
func fetchFileHttp(url string) (path string, err error) {

	// fetch the file content
	resp, err := http.Get(url)
	if err != nil {
		return path, errors.Wrapf(err, "failed to fetch file content")
	}
	defer resp.Body.Close()

	// read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return path, errors.Wrapf(err, "failed to read response content")
	}

	// create temp file
	tmpFile, err := ioutil.TempFile(os.TempDir(), "parent-devfile-*.yaml")
	if err != nil {
		return path, errors.Wrapf(err, "failed to create temporary file")
	}

	// write the content to temp file
	if _, err = tmpFile.Write(body); err != nil {
		return path, errors.Wrapf(err, "failed to write to temp file")
	}

	// Close the file
	if err := tmpFile.Close(); err != nil {
		return path, errors.Wrapf(err, "failed to close temporary file")
	}

	// Successfull
	glog.V(4).Infof("created parent devfile at '%s'", tmpFile.Name())
	return tmpFile.Name(), nil
}

// fetchFilePath validates the given uri, if the given uri is:
// file: returns the file path
// http(s) url: fetches the file and returns file path
func fetchFilePath(uri string) (path string, err error) {

	var u *url.URL

	// parse uri
	if u, err = url.Parse(uri); err != nil {
		return path, errors.Wrapf(err, "failed to parse uri '%s'", uri)
	}

	// uri should directly point to yaml file, should not have user, fragment, query
	if u.Scheme == "" || u.Path == "" || u.User != nil || u.RawPath != "" || u.ForceQuery || u.RawQuery != "" || u.Fragment != "" {
		return path, errors.Errorf("invalid devfile uri '%s'", uri)
	}

	switch u.Scheme {
	case fileScheme:
		path = u.Host + u.Path
		return path, nil
	case httpsScheme, httpScheme:
		return fetchFileHttp(uri)
	default:
		return path, errors.Errorf("protocol '%s' not supported for fetching remote devfiles", u.Scheme)
	}
}
