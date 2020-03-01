package main

import (
	"github.com/michaelhenkel/ckube/cmd"

	log "github.com/sirupsen/logrus"
)

// GlobalConfig is the global tool configuration
type GlobalConfig struct {
	Pkg PkgConfig `yaml:"pkg"`
}

// PkgConfig is the config specific to the `pkg` subcommand
type PkgConfig struct {
	// ContentTrustCommand is passed to `sh -c` and the stdout
	// (including whitespace and \n) is set as the content trust
	// passphrase. Can be used to execute a password manager.
	ContentTrustCommand string `yaml:"content-trust-passphrase-command"`
}

var (
	defaultLogFormatter = &log.TextFormatter{}

	// Config is the global tool configuration
	Config = GlobalConfig{}
)

func main() {
	cmd.Execute()
}
