package authz

import base "github.com/devpablocristo/core/saas/go/authz"

const (
	ScopeToolsRead         = base.ScopeToolsRead
	ScopeToolsWrite        = base.ScopeToolsWrite
	ScopePolicyRead        = base.ScopePolicyRead
	ScopePolicyWrite       = base.ScopePolicyWrite
	ScopeEgressRead        = base.ScopeEgressRead
	ScopeEgressWrite       = base.ScopeEgressWrite
	ScopeAuditRead         = base.ScopeAuditRead
	ScopeGatewayRun        = base.ScopeGatewayRun
	ScopeGatewaySimulate   = base.ScopeGatewaySimulate
	ScopeMCPRead           = base.ScopeMCPRead
	ScopeMCPCall           = base.ScopeMCPCall
	ScopeA2ACall           = base.ScopeA2ACall
	ScopeAdminConsoleRead  = base.ScopeAdminConsoleRead
	ScopeAdminConsoleWrite = base.ScopeAdminConsoleWrite
	ScopeSecretsAdmin      = base.ScopeSecretsAdmin
)

var (
	HasScope    = base.HasScope
	HasAnyScope = base.HasAnyScope
	IsRole      = base.IsRole
)
