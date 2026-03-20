package entitlement

import (
	"strings"

	"github.com/devpablocristo/core/saas/domain"
)

func NormalizePlan(raw string) domain.PlanCode {
	switch strings.TrimSpace(strings.ToLower(raw)) {
	case string(domain.PlanGrowth):
		return domain.PlanGrowth
	case string(domain.PlanEnterprise):
		return domain.PlanEnterprise
	default:
		return domain.PlanStarter
	}
}

func DefaultHardLimits(plan domain.PlanCode) domain.HardLimits {
	switch NormalizePlan(string(plan)) {
	case domain.PlanGrowth:
		return domain.HardLimits{
			ToolsMax:           50,
			RunRPM:             120,
			AuditRetentionDays: 90,
		}
	case domain.PlanEnterprise:
		return domain.HardLimits{
			ToolsMax:           250,
			RunRPM:             600,
			AuditRetentionDays: 365,
		}
	default:
		return domain.HardLimits{
			ToolsMax:           10,
			RunRPM:             30,
			AuditRetentionDays: 30,
		}
	}
}

func CanUse(current, limit int) bool {
	if limit <= 0 {
		return false
	}
	return current < limit
}
