package lib

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"runtime"

	log "github.com/Sirupsen/logrus"
)

var (
	errResource          = errors.New("invalid webfinger resource")
	errNotYetImplemented = errors.New("Not yet implemented")
	errTooManyRedirects  = errors.New("Too many redirects")
	errHTTPRedirect      = errors.New("Redirect to non-https server")
	errHTTPCode          = errors.New("Received unexpected http code")
	errSubjectMissmatch  = errors.New("Subject doesn't match resource")
)

// PassThru wraps an existing io.Reader.
//
// It simply forwards the Read() call, while displaying
// the results from individual calls to it.
type PassThru struct {
	io.Reader
	total    int64 // Total # of bytes transferred
	length   int64 // Expected length
	progress float64
}

func (d *Dkenv) DownloadDocker(version string) error {
	resp, err := d.getHttp(version)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	readerpt := &PassThru{Reader: resp.Body, length: resp.ContentLength}

	body, err := ioutil.ReadAll(readerpt)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(d.DkenvDir+"/docker-"+version, body, 0755); err != nil {
		return fmt.Errorf("Error(s) writing docker binary: %v", err)
	}

	return nil
}

func (d *Dkenv) getHttp(version string) (*http.Response, error) {
	client := &http.Client{
		CheckRedirect: redirectPolicyFunc,
	}

	var system string

	switch {
	case runtime.GOOS == "windows":
		system = "Windows"
	case runtime.GOOS == "linux":
		system = "Linux"
	case runtime.GOOS == "darwin":
		system = "Darwin"
	default:
		return nil, fmt.Errorf("Unsupported system type - %v", runtime.GOOS)
	}

	resp, err := client.Get("https://get.docker.com/builds/" + system + "/x86_64/docker-" + version)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("No such docker version '%v'", version)
	}

	return resp, nil

}

func (pt *PassThru) Read(p []byte) (int, error) {
	n, err := pt.Reader.Read(p)
	if n > 0 {
		pt.total += int64(n)
		percentage := float64(pt.total) / float64(pt.length) * float64(100)

		if percentage-pt.progress > 2 {
			log.Debugf("Transferred %d percent", int(percentage))
			pt.progress = percentage
		}
	}

	return n, err
}

func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	if len(via) > 10 {
		return errTooManyRedirects
	}

	if req.URL.Scheme != "https" {
		return errHTTPRedirect
	}

	return nil
}
