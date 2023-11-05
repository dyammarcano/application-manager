package main

import (
	"github.com/dyammarcano/application-manager/cmd"
	"github.com/dyammarcano/application-manager/internal/service"
)

var (
	Version    = "v0.0.1-manual-build"
	CommitHash string
	Date       string
)

func main() {
	service.Init(Version, CommitHash, Date)
	cmd.Execute()
}
