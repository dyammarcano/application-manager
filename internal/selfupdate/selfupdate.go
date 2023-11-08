package selfupdate

type (
	SelfUpdate struct {
		Url       string
		Signature string
		Found     bool
	}
)

func (s *SelfUpdate) Check() bool {
	return s.Found
}

func (s *SelfUpdate) Update() error {
	return nil
}

func (s *SelfUpdate) UpdateTo(version string) error {
	return nil
}

func (s *SelfUpdate) UpdateToLatest() error {
	return nil
}
