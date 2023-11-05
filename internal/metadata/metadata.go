package metadata

import (
	"encoding/json"
	"fmt"
)

type (
	Metadata struct {
		ApplicationVersion string   `json:"application_version"`
		CommitHash         string   `json:"commit_hash"`
		GoVersion          string   `json:"go_version"`
		ReleaseDate        string   `json:"release_date"`
		CommitTag          string   `json:"commit_tag"`
		Runtime            *Runtime `json:"runtime"`
	}

	Runtime struct {
		Arch string `json:"arch"`
		Goos string `json:"goos"`
	}
)

func (m *Metadata) String() {
	data, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(data))
}
