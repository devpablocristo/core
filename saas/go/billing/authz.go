package billing

import "github.com/devpablocristo/core/saas/go/authz"

const (
	ScopeAdminConsoleRead  = "admin:console:read"
	ScopeAdminConsoleWrite = "admin:console:write"
)

func canAccess(scopes []string, role *string, required string) bool {
	if authz.IsRole(role, "admin", "secops") {
		return true
	}
	return authz.HasScope(scopes, required)
}
