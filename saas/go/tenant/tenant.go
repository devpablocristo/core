package tenant

import (
	"strings"
	"unicode"

	"github.com/devpablocristo/core/saas/go/domain"
)

func NormalizeSlug(raw string) string {
	raw = strings.TrimSpace(strings.ToLower(raw))
	if raw == "" {
		return ""
	}

	var out []rune
	lastDash := false
	for _, r := range raw {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			out = append(out, r)
			lastDash = false
		case unicode.IsSpace(r) || r == '-' || r == '_' || r == '/':
			if !lastDash && len(out) > 0 {
				out = append(out, '-')
				lastDash = true
			}
		}
	}
	return strings.Trim(string(out), "-")
}

func NormalizeRole(raw string) string {
	raw = strings.TrimSpace(strings.ToLower(raw))
	switch raw {
	case "owner", "admin", "secops", "viewer":
		return raw
	default:
		return "viewer"
	}
}

func NewMembership(tenantID, userID, role string) domain.Membership {
	return domain.Membership{
		TenantID: tenantID,
		UserID:   userID,
		Role:     NormalizeRole(role),
	}
}
