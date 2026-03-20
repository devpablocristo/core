package domain

import kerneldomain "github.com/devpablocristo/core/governance/kernel/usecases/domain"

type Decision = kerneldomain.Decision
type Evaluation = kerneldomain.Evaluation

const (
	DecisionAllow           = kerneldomain.DecisionAllow
	DecisionDeny            = kerneldomain.DecisionDeny
	DecisionRequireApproval = kerneldomain.DecisionRequireApproval
)
