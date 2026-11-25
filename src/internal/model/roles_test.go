package model

import "testing"

func TestIsDefaultRoleName(t *testing.T) {
	for _, name := range []string{"Administrator", "Operator", "Viewer"} {
		if !IsDefaultRoleName(name) {
			t.Errorf("IsDefaultRoleName(%q) = false, want true", name)
		}
	}
	for _, name := range []string{"Support", "administrator"} {
		if IsDefaultRoleName(name) {
			t.Errorf("IsDefaultRoleName(%q) = true, want false", name)
		}
	}
}
