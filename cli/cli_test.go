package cli

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
)

var (
	testBinDir   = "/usr/bin"
	testHomeDir  = "/tmp/test_home_dir"
	testDkenvDir = "" // Get's set via init()
)

func init() {
	// setup dummy dirs
	tmpDir, err := ioutil.TempDir("", ".dkenv")
	if err != nil {
		fmt.Println("Unable to create tmp dir for testing: %v", err)
		os.Exit(1)
	}

	testDkenvDir = tmpDir
	defer cleanUp(testDkenvDir)
}

func TestNew(t *testing.T) {
	c := New(&testBinDir, &testHomeDir, &testDkenvDir)

	assert.Equal(t, testBinDir, *c.BinDir)
	assert.Equal(t, testHomeDir, *c.HomeDir)
	assert.Equal(t, testDkenvDir, *c.DkenvDir)
}

func TestSetHomeDir(t *testing.T) {
	// Blank homedir - acquire it automatically (happy path)
	homeDir, homeDirErr := homedir.Dir()
	assert.NoError(t, homeDirErr)

	var blankDir string

	c1 := New(&testBinDir, &blankDir, &testDkenvDir)
	err1 := c1.setHomeDir()

	assert.NoError(t, err1)
	assert.Equal(t, homeDir, *c1.HomeDir)

	// Dir test with not-existent dir
	badDir := fmt.Sprintf("/tmp/bad_dir_%v", rand.Intn(99999))
	c2 := New(&testBinDir, &badDir, &testDkenvDir)
	err2 := c2.setHomeDir()

	assert.Error(t, err2)

	// Verify right trim takes place
	homeDirWithSlash := homeDir + "/"

	c3 := New(&testBinDir, &homeDirWithSlash, &testDkenvDir)
	err3 := c3.setHomeDir()

	assert.NoError(t, err3)
	assert.Equal(t, homeDir, *c3.HomeDir)
}

func TestHandleBinDir(t *testing.T) {
	// Non-existent dir check
	badDir := fmt.Sprintf("/tmp/bad_dir_%v", rand.Intn(99999))
	c1 := New(&badDir, &testHomeDir, &testDkenvDir)
	err1 := c1.handleBinDir()

	assert.Error(t, err1)

	// strip trailing slash check
	testBinDirWithSlash := testBinDir + "/"
	c2 := New(&testBinDirWithSlash, &testHomeDir, &testDkenvDir)
	err2 := c2.handleBinDir()

	assert.NoError(t, err2)
	assert.Equal(t, testBinDir, *c2.BinDir)
}

func TestHandleDkenvDir(t *testing.T) {
	// Create a tmp home dir
	tmpHomeDir, err := ioutil.TempDir("", "dkenv_temp_home_dir")
	if err != nil {
		t.Fatalf("Unable to create temp dir for testing: %v", err)
	}
	defer cleanUp(tmpHomeDir)

	// Happy path + auto create dkenv dir
	dkenvDir1 := "~/.dkenv"
	c1 := New(&testBinDir, &tmpHomeDir, &dkenvDir1)
	err1 := c1.handleDkenvDir()

	assert.NoError(t, err1)
	assert.Equal(t, tmpHomeDir+"/.dkenv", *c1.DkenvDir)

	// Assert directory has been created
	if _, err := os.Stat(*c1.DkenvDir); os.IsNotExist(err) {
		t.Errorf("%v not created after handleDkenvDir()", *c1.DkenvDir)
	}

	// Trailing slash check
	testDkenvDirWithSlash := testDkenvDir + "/"
	c2 := New(&testBinDir, &testHomeDir, &testDkenvDirWithSlash)
	err2 := c2.handleDkenvDir()

	assert.NoError(t, err2)
	assert.Equal(t, testDkenvDir, *c2.DkenvDir)

	// Test inability to create a dkenv dir
	tmpFile, err := ioutil.TempFile("", "dkenv_temp_file")
	if err != nil {
		t.Fatalf("Unable to create a temp file during testing: %v", err)
	}
	defer cleanUp(tmpFile.Name())

	tmpFilename := tmpFile.Name()

	c3 := New(&testBinDir, &testHomeDir, &tmpFilename)
	err3 := c3.handleDkenvDir()
	assert.Error(t, err3)
	assert.Contains(t, err3.Error(), "Error creating dkenv dir")
}

func TestHandleArgs(t *testing.T) {
	homeDir, homeDirErr := homedir.Dir()
	assert.NoError(t, homeDirErr)

	// bad homedir should return error
	badDir := fmt.Sprintf("/tmp/bad_dir_%v", rand.Intn(99999))

	c1 := New(&testBinDir, &badDir, &testDkenvDir)
	assert.Error(t, c1.HandleArgs())

	// bad bindir should return error
	c2 := New(&badDir, &homeDir, &testDkenvDir)
	assert.Error(t, c2.HandleArgs())

	// pre-existing tmp file for dkenv
	tmpFile, err := ioutil.TempFile("", "dkenv_temp_file")
	if err != nil {
		t.Fatalf("Unable to create a temp file during testing: %v", err)
	}
	defer cleanUp(tmpFile.Name())

	tmpFilename := tmpFile.Name()

	c3 := New(&testBinDir, &homeDir, &tmpFilename)
	err3 := c3.HandleArgs()
	assert.Error(t, err3)
	assert.Contains(t, err3.Error(), "Error creating dkenv")

	// happy path
	c4 := New(&testBinDir, &homeDir, &testDkenvDir)
	assert.NoError(t, c4.HandleArgs())
}

func TestIsDirAndExists(t *testing.T) {
	// Test non-existent dir
	badDir := fmt.Sprintf("/tmp/bad_dir_%v", rand.Intn(99999))

	test1, err1 := IsDirAndExists(badDir)
	assert.NoError(t, err1)
	assert.False(t, test1)

	// Test w/ file (instead of dir)
	tmpFile, err2 := ioutil.TempFile("", "dkenv_temp_file")
	if err2 != nil {
		t.Fatalf("Unable to create a temp file during testing: %v", err2)
	}
	defer cleanUp(tmpFile.Name())

	test2, err2 := IsDirAndExists(tmpFile.Name())
	assert.NoError(t, err2, "Should get an error when given a file instead of dir")
	assert.False(t, test2)

	// Happy path
	tmpDir, err3 := ioutil.TempDir("", "dkenv_temp_dir")
	if err3 != nil {
		t.Fatalf("Unable to create temp dir for testing: %v", err3)
	}
	defer cleanUp(tmpDir)

	test3, err3 := IsDirAndExists(tmpDir)
	assert.NoError(t, err3)
	assert.True(t, test3)
}

func cleanUp(tmpDir string) {
	if err := os.RemoveAll(tmpDir); err != nil {
		fmt.Println("Unable to clean up tmp dir: %v", err)
		os.Exit(1)
	}
}
