package notifiers

import (
	"errors"
	"testing"
)

func TestSendWithRetry(t *testing.T) {
	attempts := 0
	err := SendWithRetry(func() error {
		attempts++
		if attempts == 1 {
			return errors.New("temporary failure")
		}
		return nil
	}, true)
	if err != nil || attempts != 2 {
		t.Fatalf("SendWithRetry() = %v after %d attempts, want success after 2", err, attempts)
	}

	attempts = 0
	err = SendWithRetry(func() error {
		attempts++
		return errors.New("failure")
	}, false)
	if err == nil || attempts != 1 {
		t.Fatalf("SendWithRetry(retry=false) = %v after %d attempts, want one failed attempt", err, attempts)
	}
}
