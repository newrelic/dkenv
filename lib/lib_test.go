package lib

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	// Set during init()
	testDkenvDir = ""
	testBinDir   = ""
)

// Create temp dirs necessary for tests
func init() {
	tmpBinDir, err := ioutil.TempDir("", "tmp_dkenv_bin")
	if err != nil {
		fmt.Println("Unable to create tmp dir for testing: %v", err)
		os.Exit(1)
	}

	testBinDir = tmpBinDir
	defer cleanUp(testDkenvDir)

	tmpDkenvDir, err := ioutil.TempDir("", "tmp_dkenv_dir")
	if err != nil {
		fmt.Println("Unable to create tmp dir for testing: %v", err)
		os.Exit(1)
	}

	testDkenvDir = tmpDkenvDir
	defer cleanUp(testDkenvDir)

}

func TestNew(t *testing.T) {
	c := New(testDkenvDir, testBinDir)
	assert.IsType(t, &Dkenv{}, c)
	assert.Equal(t, testDkenvDir, c.DkenvDir)
	assert.Equal(t, testBinDir, c.BinDir)
}

func TestListAction(t *testing.T) {

}

func TestListInstalled(t *testing.T) {

}

func TestIsInstalled(t *testing.T) {

}

func TestFetchVersionAction(t *testing.T) {

}

func TestApiToVersion(t *testing.T) {

}

func TestUpdateSymlink(t *testing.T) {

}

func cleanUp(tmpDir string) {
	if err := os.RemoveAll(tmpDir); err != nil {
		fmt.Println("Unable to clean up tmp dir: %v", err)
		os.Exit(1)
	}
}
