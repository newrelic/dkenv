package lib

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"

	log "github.com/Sirupsen/logrus"
)

var (
	apiVersions = map[string]string{
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
		"1.22": "1.10.1",
	}
)

type Dkenv struct {
	DkenvDir string
	BinDir   string
}

func New(dkenvDir, binDir string) *Dkenv {
	return &Dkenv{
		DkenvDir: dkenvDir,
		BinDir:   binDir,
	}
}

func (d *Dkenv) ListAction() error {
	installed, err := d.listInstalled()
	if err != nil {
		return err
	}

	if len(installed) == 0 {
		log.Warning("No installed Docker binaries found!")
		return nil
	} else {
		log.Infof("Found %v installed docker binaries", len(installed))
		log.Info("") // blank line, for the pretty
	}

	for i := 0; i < len(installed); i++ {
		log.Infof("%v: %v", i+1, installed[i])
	}

	return nil
}

func (d *Dkenv) FetchVersionAction(version string, api bool) error {
	clientVersion := version

	if api {
		var err error

		clientVersion, err = ApiToVersion(version)
		if err != nil {
			return err
		}

		log.Infof("Found client '%v' for API version '%v'", clientVersion, version)
	}

	if !d.isInstalled(clientVersion) {
		// Attempt to download it
		log.Infof("Docker version %v not found - attempting to download...", clientVersion)

		if err := d.DownloadDocker(clientVersion); err != nil {
			return err
		}
	} else {
		log.Infof("Docker version %v already installed!", clientVersion)
	}

	// Update symlink
	if err := d.UpdateSymlink(clientVersion); err != nil {
		return err
	}

	return nil
}

// Return a list of found docker binaries
func (d *Dkenv) listInstalled() ([]string, error) {
	fileList, err := ioutil.ReadDir(d.DkenvDir)
	if err != nil {
		return nil, err
	}

	found := make([]string, 0)

	for _, filename := range fileList {
		match, _ := regexp.MatchString(`docker-.+`, filename.Name())
		if match {
			found = append(found, filename.Name())
		}
	}

	return found, nil
}

func (d *Dkenv) isInstalled(version string) bool {
	if _, err := os.Stat(d.DkenvDir + "/docker-" + version); err == nil {
		return true
	}

	return false
}

// Symlink docker to a dkenv docker binary
//
// - If destination exists AND is *not* a symlink:
//     - backup destination file
//     - create symlink
// - If destination exists AND *is* a symlink:
//     - check if the symlink already matches destination (noop if it does)
//     - remove and overwrite if symlink is incorrect
// - If destination does not exist:
//     - create symlink
//
// TODO This will be a pain to test - refactor when possible.
func (d *Dkenv) UpdateSymlink(version string) error {
	src := d.DkenvDir + "/docker-" + version
	dst := d.BinDir + "/docker"

	// Verify that the src file actually exists
	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("Unable to lookup source binary '%v': %v", src, err)
	}

	// Check if file exists; bail if real error
	_, err := os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("Problems verifying existing docker symlink: %v", err)
	}

	// File does not exist, let's create it
	if err != nil && os.IsNotExist(err) {
		log.Infof("Creating symlink for '%v' -> '%v'", dst, src)

		if err := os.Symlink(src, dst); err != nil {
			return fmt.Errorf("Unable to create new docker symlink: %v", err)
		}

		return nil
	}

	// File exists; check if symlink
	fi, err := os.Lstat(dst)
	if err != nil {
		return fmt.Errorf("Unable to stat existing docker file '%v': %v", dst, err)
	}

	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		tmpDst, err := os.Readlink(dst)
		if err != nil {
			return fmt.Errorf("Unable to lookup link info for existing docker symlink: %v", err)
		}

		if tmpDst == src {
			log.Infof("'%v' already pointing to '%v' - nothing to do!", dst, src)
			return nil
		}

		if err := os.Remove(dst); err != nil {
			return fmt.Errorf("Unable to remove old symlink: %v", err)
		}
	} else {
		backupName := dst + fmt.Sprintf(".%v.dkenv", rand.Intn(9999))

		log.Infof("Backing up existing docker file %v to %v", dst, backupName)

		if err := os.Rename(dst, backupName); err != nil {
			return fmt.Errorf("Unable to backup existing docker file: %v", err)
		}
	}

	// Create/overwrite old symlink
	log.Infof("Creating symlink for '%v' -> '%v'", dst, src)

	if err := os.Symlink(src, dst); err != nil {
		return fmt.Errorf("Unable to create new docker symlink: %v", err)
	}

	return nil
}

func ApiToVersion(version string) (string, error) {
	if val, ok := apiVersions[version]; ok {
		return val, nil
	}

	return "", errors.New("Invalid API Version")
}
