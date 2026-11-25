package services

import "testing"

func TestSecurityPolicyParsingAndEnforcement(t *testing.T) {
	policy, err := ParseSecurityPolicy(`{"ip_allowlist":["192.0.2.0/24"],"api_key_allowlist":[7,9]}`)
	if err != nil {
		t.Fatal(err)
	}
	if !policy.Allows("192.0.2.12", 9) {
		t.Fatal("expected listed IP and API key to be allowed")
	}
	for _, test := range []struct {
		ip  string
		key int64
	}{
		{"192.0.3.12", 9},
		{"192.0.2.12", 8},
		{"not-an-ip", 9},
	} {
		if policy.Allows(test.ip, test.key) {
			t.Fatalf("unexpectedly allowed IP %q and API key %d", test.ip, test.key)
		}
	}

	open, err := ParseSecurityPolicy(`{}`)
	if err != nil || !open.Allows("198.51.100.4", 123) {
		t.Fatalf("empty policy must preserve existing access: policy=%+v err=%v", open, err)
	}
}

func TestSecurityPolicyRejectsMalformedConfiguration(t *testing.T) {
	for _, value := range []string{
		`{"ip_allowlist":["not-an-ip"]}`,
		`{"api_key_allowlist":[0]}`,
		`{"unknown":true}`,
		`{} {}`,
		`null`,
	} {
		if _, err := ParseSecurityPolicy(value); err == nil {
			t.Errorf("ParseSecurityPolicy(%q) unexpectedly succeeded", value)
		}
	}
}
