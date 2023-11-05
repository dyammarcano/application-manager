package main

import (
	"github.com/dyammarcano/application-manager/cmd"
	"github.com/dyammarcano/application-manager/internal/application"
)

var (
	Version    = "v0.0.1-manual-build"
	CommitHash string
	Date       string
)

func main() {
	manager := application.NewApplicationManager(Version, CommitHash, Date)
	cmd.Execute(manager)
}
