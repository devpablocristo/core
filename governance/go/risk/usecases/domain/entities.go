package domain

import kerneldomain "github.com/devpablocristo/core/governance/go/kernel/usecases/domain"

type RiskLevel = kerneldomain.RiskLevel
type RiskFactor = kerneldomain.RiskFactor
type RiskAssessment = kerneldomain.RiskAssessment

const (
	RiskLow    = kerneldomain.RiskLow
	RiskMedium = kerneldomain.RiskMedium
	RiskHigh   = kerneldomain.RiskHigh
)
