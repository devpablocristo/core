package domain

import kerneldomain "github.com/devpablocristo/core/governance/go/kernel/usecases/domain"

type Policy = kerneldomain.Policy
type PolicyMode = kerneldomain.PolicyMode

const (
	PolicyModeEnforce = kerneldomain.PolicyModeEnforce
	PolicyModeShadow  = kerneldomain.PolicyModeShadow
)
