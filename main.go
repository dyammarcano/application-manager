package main

import (
	"github.com/dyammarcano/template-go/cmd"
	"github.com/dyammarcano/template-go/internal/application"
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
