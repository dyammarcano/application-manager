package selfupdate

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type (
	SelfUpdate struct {
		ApiURL         string
		BinURL         string
		Dir            string
		Signature      string
		CurrentVersion string
		ForceCheck     bool
		CmdName        string
		EnableLogging  bool
		NoConfirm      bool
		NoProgress     bool
		NoRestart      bool
		noUmask        bool
		root           bool
		disableCurl    bool
		disableWget    bool
		report         chan any
		wg             sync.WaitGroup
		ticker         *time.Ticker
		checkEvery     time.Duration
		nextCheck      time.Time
		status         string
	}
)

func NewSelfUpdate() *SelfUpdate {
	return &SelfUpdate{
		ApiURL:         "",
		BinURL:         "",
		Dir:            "",
		Signature:      "",
		CurrentVersion: "",
		ForceCheck:     false,
		CmdName:        "",
		EnableLogging:  false,
		NoConfirm:      false,
		NoProgress:     false,
		NoRestart:      false,
		noUmask:        false,
		root:           false,
		disableCurl:    false,
		disableWget:    false,
		report:         make(chan any),
		wg:             sync.WaitGroup{},
		ticker:         nil,
		checkEvery:     24 * time.Hour,
		nextCheck:      time.Now(),
	}
}

func (s *SelfUpdate) Check() bool {
	return false
}

func (s *SelfUpdate) UpdateBackgroundRun() {
	s.wg.Add(1)
	s.status = "Checking for updates..."

	go func() {
		defer s.wg.Done()
		s.ticker = time.NewTicker(s.checkEvery)
		if err := s.UpdateToLatest(); err != nil {
			s.report <- err
		}
	}()
}

func (s *SelfUpdate) UpdateTo(version string) error {
	return nil
}

func (s *SelfUpdate) UpdateToLatest() error {
	if !s.Check() {
		s.status = "No updates found, next check in 24 hours."
		s.nextCheck = time.Now().Add(s.checkEvery)
		return nil
	}

	s.downloadFile()
	return nil
}

func (s *SelfUpdate) getExecRelativeDir(dir string) string {
	filename, _ := os.Executable()
	path := filepath.Join(filepath.Dir(filename), dir)
	return path
}

func (s *SelfUpdate) Status() string {
	return s.status
}

func (s *SelfUpdate) Stop() {
	s.ticker.Stop()
	s.wg.Wait()
}

func (s *SelfUpdate) downloadFile() {
	s.status = "Update found, verify checksum..."
	if !s.verify() {
		s.status = "Update failed, checksum mismatch..."
		return
	}
	s.status = "Updating..."
	s.doReplace()
}

func (s *SelfUpdate) verify() bool {
	//todo: unzip file
	//todo: verify checksum

	return true
}

func (s *SelfUpdate) doReplace() {
	//todo copy file to new location
	// change permission
	// restart
}

func canUpdate() (err error) {
	// get the directory the file exists in
	path, err := os.Executable()
	if err != nil {
		return
	}

	fileDir := filepath.Dir(path)
	fileName := filepath.Base(path)

	// attempt to open a file in the file's directory
	newPath := filepath.Join(fileDir, fmt.Sprintf(".%s.new", fileName))
	fp, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return
	}
	fp.Close()

	_ = os.Remove(newPath)
	return
}
