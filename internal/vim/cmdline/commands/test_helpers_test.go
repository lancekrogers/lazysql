package commands

type fakeStatus struct {
	infos  []string
	errors []string
}

func (f *fakeStatus) Info(message string) {
	f.infos = append(f.infos, message)
}

func (f *fakeStatus) Error(message string) {
	f.errors = append(f.errors, message)
}

type fakeApp struct {
	stopped bool
}

func (f *fakeApp) Stop() {
	f.stopped = true
}
