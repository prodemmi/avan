package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/netip"
	"strings"
)

const APISecurityPolicySettingKey = "api-security-policy"

type SecurityPolicy struct {
	IPAllowlist     []string `json:"ip_allowlist"`
	APIKeyAllowlist []int64  `json:"api_key_allowlist"`
}

func ParseSecurityPolicy(value string) (SecurityPolicy, error) {
	var policy SecurityPolicy
	value = strings.TrimSpace(value)
	if value == "" {
		return policy, nil
	}
	if !strings.HasPrefix(value, "{") {
		return policy, fmt.Errorf("security policy must be a JSON object")
	}

	decoder := json.NewDecoder(strings.NewReader(value))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&policy); err != nil {
		return policy, fmt.Errorf("invalid security policy: %w", err)
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return policy, fmt.Errorf("security policy must contain one JSON object")
	}

	for _, entry := range policy.IPAllowlist {
		if _, err := securityPrefix(entry); err != nil {
			return policy, fmt.Errorf("invalid IP allowlist entry %q: %w", entry, err)
		}
	}
	for _, id := range policy.APIKeyAllowlist {
		if id <= 0 {
			return policy, fmt.Errorf("API key IDs must be positive")
		}
	}
	return policy, nil
}

func (p SecurityPolicy) Allows(clientIP string, apiKeyID int64) bool {
	if len(p.IPAllowlist) > 0 {
		ip, err := netip.ParseAddr(clientIP)
		if err != nil {
			return false
		}
		allowed := false
		for _, entry := range p.IPAllowlist {
			prefix, err := securityPrefix(entry)
			if err == nil && prefix.Contains(ip) {
				allowed = true
				break
			}
		}
		if !allowed {
			return false
		}
	}
	if len(p.APIKeyAllowlist) > 0 {
		for _, id := range p.APIKeyAllowlist {
			if id == apiKeyID {
				return true
			}
		}
		return false
	}
	return true
}

func securityPrefix(value string) (netip.Prefix, error) {
	if prefix, err := netip.ParsePrefix(value); err == nil {
		if prefix.Addr().Zone() != "" {
			return netip.Prefix{}, fmt.Errorf("scoped IP addresses are not supported")
		}
		return prefix.Masked(), nil
	}
	ip, err := netip.ParseAddr(value)
	if err != nil {
		return netip.Prefix{}, err
	}
	if ip.Zone() != "" {
		return netip.Prefix{}, fmt.Errorf("scoped IP addresses are not supported")
	}
	return netip.PrefixFrom(ip, ip.BitLen()), nil
}
