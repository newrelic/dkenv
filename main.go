package main

import (
	"github.com/newrelic/dkenv/cli"
	"github.com/newrelic/dkenv/lib"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	VERSION           = "1.1.0"
	DEFAULT_BINDIR    = "/usr/local/bin"
	DEFAULT_DKENV_DIR = "~/.dkenv"
)

var (
	binDir   = kingpin.Flag("bindir", "Directory to create symlinks for Docker binaries").Default(DEFAULT_BINDIR).String()
	homeDir  = kingpin.Flag("homedir", "Override automatically found homedir").String()
	dkenvDir = kingpin.Flag("dkenvdir", "Directory to store Docker binaries").Default(DEFAULT_DKENV_DIR).String()
	debug    = kingpin.Flag("debug", "Enable debug output").Short('d').Bool()

	// Commands
	client    = kingpin.Command("client", "Download/switch Docker binary by *client* version")
	clientArg = client.Arg("version", "Docker client version").Required().String()
	api       = kingpin.Command("api", "Download/switch Docker binary by *API* version")
	apiArg    = api.Arg("version", "Docker API version").Required().String()

	// Not sure if there is any other, non-hacky way to set a command with no
	// args and retain the ability to see if it was called.
	list = kingpin.Command("list", "List downloaded/existing Docker binaries").Action(func(c *kingpin.ParseContext) error {
		listAction = true
		return nil
	})

	listAction = false
)

func init() {
	kingpin.Version(VERSION)
	kingpin.Parse()

	c := cli.New(binDir, homeDir, dkenvDir)
	if err := c.HandleArgs(); err != nil {
		log.Fatalf("Argument error: %v", err)
	}

	if *debug {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	d := lib.New(*dkenvDir, *binDir)

	var err error

	switch {
	case listAction:
		err = d.ListAction()
	case *apiArg != "":
		err = d.FetchVersionAction(*apiArg, true)
	case *clientArg != "":
		err = d.FetchVersionAction(*clientArg, false)
	}

	if err != nil {
		log.Fatalf(err.Error())
	}
}
