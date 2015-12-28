package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
)

type Cli struct {
	BinDir   *string
	HomeDir  *string
	DkenvDir *string
}

func New(binDir, homeDir, dkenvDir *string) *Cli {
	c := &Cli{
		BinDir:   binDir,
		HomeDir:  homeDir,
		DkenvDir: dkenvDir,
	}

	return c
}

// Wrapper for each specific arg handler
func (c *Cli) HandleArgs() error {
	if err := c.setHomeDir(); err != nil {
		return err
	}

	if err := c.handleBinDir(); err != nil {
		return err
	}

	if err := c.handleDkenvDir(); err != nil {
		return err
	}

	return nil
}

// If homedir is not set, set it; if it doesn't exist, bail out
func (c *Cli) setHomeDir() error {
	if *c.HomeDir == "" {
		homeDir, err := homedir.Dir()
		if err != nil {
			return fmt.Errorf("Unable to fetch current home dir: %v", err)
		}

		*c.HomeDir = homeDir
	}

	// Ensure it exists
	isDir, err := IsDirAndExists(*c.HomeDir)
	if err != nil {
		return fmt.Errorf("Home directory error: %v", err)
	}

	if !isDir {
		return errors.New("Home directory does not exist!")
	}

	// Strip trailing slash
	*c.HomeDir = strings.TrimRight(*c.HomeDir, "/")

	return nil
}

func (c *Cli) handleBinDir() error {
	// Ensure it exists
	isDir, err := IsDirAndExists(*c.BinDir)
	if err != nil {
		return fmt.Errorf("Bin directory error: %v", err)
	}

	if !isDir {
		return errors.New("Bin directory does not exist!")
	}

	// Strip trailing slash
	*c.BinDir = strings.TrimRight(*c.BinDir, "/")

	return nil
}

// Check if dkenv dir exists; if not, create it
func (c *Cli) handleDkenvDir() error {
	// Expand the dir
	if *c.DkenvDir == "~/.dkenv" {
		*c.DkenvDir = *c.HomeDir + "/.dkenv"
	}

	isDir, err := IsDirAndExists(*c.DkenvDir)
	if err != nil {
		return fmt.Errorf("Dkenv directory error: %v", err)
	}

	if !isDir {
		// Create the dir
		if err := os.Mkdir(*c.DkenvDir, 0700); err != nil {
			return fmt.Errorf("Error creating dkenv dir (%v): %v", *c.DkenvDir, err)
		}
	}

	// Strip trailing slash
	*c.DkenvDir = strings.TrimRight(*c.DkenvDir, "/")

	return nil
}

func IsDirAndExists(dir string) (bool, error) {
	fi, err := os.Stat(dir)
	if err != nil && os.IsNotExist(err) {
		return false, nil
	}

	// Regular error
	if err != nil {
		return false, err
	}

	if !fi.IsDir() {
		return false, nil
	}

	return true, nil
}
