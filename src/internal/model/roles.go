package model

// Default role names are the only system-role marker in the current schema.
func IsDefaultRoleName(name string) bool {
	switch name {
	case "Administrator", "Operator", "Viewer":
		return true
	default:
		return false
	}
}
