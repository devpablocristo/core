package ctxkeys

import base "github.com/devpablocristo/core/saas/go/ctxkeys"

type Key = base.Key

const (
	RequestID  = base.RequestID
	OrgID      = base.OrgID
	TenantID   = base.TenantID
	Actor      = base.Actor
	Role       = base.Role
	Scopes     = base.Scopes
	AuthMethod = base.AuthMethod
)
