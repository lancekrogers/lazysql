package leader

import (
	"testing"
	"time"
)

func TestTimeoutManagerStart(t *testing.T) {
	triggered := make(chan struct{})
	manager := NewTimeoutManager(10*time.Millisecond, func() {
		close(triggered)
	})

	manager.Start()
	select {
	case <-triggered:
		// ok
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout did not trigger")
	}
}

func TestTimeoutManagerReset(t *testing.T) {
	triggered := make(chan struct{})
	manager := NewTimeoutManager(30*time.Millisecond, func() {
		close(triggered)
	})

	manager.Start()
	time.Sleep(15 * time.Millisecond)
	manager.Reset()

	select {
	case <-triggered:
		t.Fatal("timeout should not have triggered after reset")
	case <-time.After(20 * time.Millisecond):
		// ok
	}

	select {
	case <-triggered:
		// ok
	case <-time.After(60 * time.Millisecond):
		t.Fatal("timeout did not trigger after reset")
	}
}

func TestTimeoutManagerStop(t *testing.T) {
	triggered := make(chan struct{})
	manager := NewTimeoutManager(20*time.Millisecond, func() {
		close(triggered)
	})

	manager.Start()
	manager.Stop()

	select {
	case <-triggered:
		t.Fatal("timeout should not trigger after stop")
	case <-time.After(40 * time.Millisecond):
		// ok
	}
}

func TestTimeoutManagerSetDuration(t *testing.T) {
	triggered := make(chan struct{})
	manager := NewTimeoutManager(50*time.Millisecond, func() {
		close(triggered)
	})

	manager.SetDuration(5 * time.Millisecond)
	manager.Start()

	select {
	case <-triggered:
		// ok
	case <-time.After(50 * time.Millisecond):
		t.Fatal("timeout did not respect updated duration")
	}
}
