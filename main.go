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

func init() {
	application.Init(Version, CommitHash, Date)
}

func main() {
	cmd.Execute()
}
