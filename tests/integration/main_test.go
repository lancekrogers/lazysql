//go:build integration
// +build integration

package integration

import (
	"os"
	"testing"
)

var sharedContainer *TestContainer

func TestMain(m *testing.M) {
	var err error
	sharedContainer, err = NewSharedContainer()
	if err != nil {
		os.Stderr.WriteString("failed to create shared container: " + err.Error() + "\n")
		os.Exit(1)
	}

	code := m.Run()

	sharedContainer.Cleanup()
	os.Exit(code)
}

func GetSharedContainer(t *testing.T) *TestContainer {
	t.Helper()
	if sharedContainer == nil {
		t.Fatal("shared container not initialized")
	}

	if err := sharedContainer.Reset(); err != nil {
		t.Fatalf("failed to reset container: %v", err)
	}

	return &TestContainer{
		container: sharedContainer.container,
		ctx:       sharedContainer.ctx,
		t:         t,
		host:      sharedContainer.host,
		port:      sharedContainer.port,
		user:      sharedContainer.user,
		password:  sharedContainer.password,
		database:  sharedContainer.database,
	}
}
