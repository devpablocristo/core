package ctxkeys

// Key tipa las context keys para evitar colisiones.
type Key string

const (
	RequestID  Key = "request_id"
	OrgID      Key = "org_id"
	TenantID   Key = "tenant_id"
	Actor      Key = "actor"
	Role       Key = "role"
	Scopes     Key = "scopes"
	AuthMethod Key = "auth_method"
)
