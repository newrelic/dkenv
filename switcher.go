package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"
)

func apiToVersion(apiVersion string) (string, error) {
	versions := map[string]string{
		"1.12": "1.0.1",
		"1.13": "1.1.2",
		"1.14": "1.2.0",
		"1.15": "1.3.3",
		"1.16": "1.4.1",
		"1.17": "1.5.0",
		"1.18": "1.6.0",
		"1.19": "1.7.1",
		"1.20": "1.8.3",
		"1.21": "1.9.1",
	}
	if val, ok := versions[apiVersion]; ok {
		return val, nil
	} else {
		return "", errors.New("Invalid API Version")
	}
}

func listDownloadedVersions() {
	usr, err := user.Current()

	if err != nil {
		log.Fatal(err)
	}

	fileList, err := ioutil.ReadDir(usr.HomeDir + "/.dkenv")

	for _, file := range fileList {
		fmt.Println(file)
	}
}

func versionDownloaded(version string) bool {
	usr, err := user.Current()

	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(usr.HomeDir + "/.dkenv/docker-" + version); err == nil {
		return true
	}
	return false
}

func createLocalLink(binDir string) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	if strings.HasSuffix(binDir, "/") {
		binDir = binDir[:len(binDir)-len("/")]
	}

	usrlocalbindocker := binDir + "/docker"

	if _, err := os.Stat(usrlocalbindocker); err == nil {
		ulbdockerFile, _ := os.Lstat(usrlocalbindocker)
		// If the file exists and is a symlink do nothing
		if ulbdockerFile.Mode()&os.ModeSymlink == os.ModeSymlink {
		} else {
			// If the file exists and is a binary, rename it.
			os.Rename(usrlocalbindocker, usrlocalbindocker+"-predkenv")
			fmt.Println("Moving " + usrlocalbindocker + " to " + usrlocalbindocker + ".predkenv")
		}
	}

	// If the file does not exist, go ahead and create a symlink
	if _, err := os.Stat(usrlocalbindocker); os.IsNotExist(err) {
		os.Symlink(usr.HomeDir+"/.dkenv/docker", usrlocalbindocker)
		fmt.Println("Creating symlink from " + usr.HomeDir + "/.dkenv/docker to /usr/local/bin/docker")
	}

}

func switchVersion(version string, binDir string) bool {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	if versionDownloaded(version) {
		os.Remove(usr.HomeDir + "/.dkenv/docker")
		os.Symlink(usr.HomeDir+"/.dkenv/docker-"+version, usr.HomeDir+"/.dkenv/docker")
		createLocalLink(binDir)

	}

	return false

}
