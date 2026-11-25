package router

import (
	"avan/internal/db/sqlc"
	"testing"
)

func TestWorkspaceAllowed(t *testing.T) {
	workspaces := []sqlc.Workspace{{ID: 3}, {ID: 8}}
	for _, test := range []struct {
		id      int64
		allowed bool
	}{
		{id: 3, allowed: true},
		{id: 8, allowed: true},
		{id: 9, allowed: false},
		{id: 0, allowed: false},
	} {
		if got := workspaceAllowed(workspaces, test.id); got != test.allowed {
			t.Errorf("workspaceAllowed(%d) = %t, want %t", test.id, got, test.allowed)
		}
	}
}
