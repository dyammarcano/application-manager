package main

import (
	"github.com/dyammarcano/application-manager/cmd"
)

var (
	Version    = "v0.0.1-manual-build"
	CommitHash string
	Date       string
)

func main() {
	cmd.Execute(Version, CommitHash, Date)
}
