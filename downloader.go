package main

import (
	"errors"
	"github.com/cheggaaa/pb"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"runtime"
)

var (
	Err_Resource          = errors.New("invalid webfinger resource")
	Err_NotYetImplemented = errors.New("Not yet implemented")
	Err_Too_Many_Redirect = errors.New("Too many redirects")
	Err_HTTP_Redirect     = errors.New("Redirect to non-https server")
	Err_HTTP_Code         = errors.New("Received unexpected http code")
	Err_Subject_Missmatch = errors.New("Subject doesn't match resource")
)

func createDotDKEnvDirectory() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	_err := os.Mkdir(usr.HomeDir+"/.dkenv/", 0700)
	if _err != nil {
	}
}

func createVersionFile(version string) (out *os.File) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	out, _err := os.Create(usr.HomeDir + "/.dkenv/docker-" + version)

	if _err != nil {
		log.Fatal(_err)
	}
	return out
}

func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	if len(via) > 10 {
		return Err_Too_Many_Redirect
	}

	if req.URL.Scheme != "https" {
		return Err_HTTP_Redirect
	}

	return nil
}

func getHttp(version string) *http.Response {
	client := &http.Client{
		CheckRedirect: redirectPolicyFunc,
	}

	system := "Darwin"

	if runtime.GOOS == "windows" {
		system = "Windows"
	}

	if runtime.GOOS == "linux" {
		system = "Linux"
	}

	if runtime.GOOS == "darwin" {
		system = "Darwin"
	}

	resp, _err := client.Get("https://get.docker.com/builds/" + system + "/x86_64/docker-" + version)

	if _err != nil {
		log.Fatal(_err)
	}

	return resp

}

func GetDocker(version string, binDir string) {
	createDotDKEnvDirectory()
	// Create the docker binary file
	// out := createVersionFile(version)

	// defer out.Close()
	// Do the htp get
	resp := getHttp(version)

	defer resp.Body.Close()

	bar := pb.New(int(resp.ContentLength)).SetUnits(pb.U_BYTES)
	bar.Start()

	rd := bar.NewProxyReader(resp.Body)

	body, _err := ioutil.ReadAll(rd)

	if _err != nil {
		log.Fatal(_err)
	}

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	_err = ioutil.WriteFile(usr.HomeDir+"/.dkenv/docker-"+version, body, 0777)

	SwitchVersion(version, binDir)

	if _err != nil {
		log.Fatal(_err)
	}

}
