package domainerr

import base "github.com/devpablocristo/core/saas/go/domainerr"

type Kind = base.Kind
type Error = base.Error

const (
	KindUnauthorized = base.KindUnauthorized
	KindForbidden    = base.KindForbidden
	KindNotFound     = base.KindNotFound
	KindValidation   = base.KindValidation
	KindConflict     = base.KindConflict
	KindInternal     = base.KindInternal
)

var (
	New          = base.New
	Newf         = base.Newf
	Unauthorized = base.Unauthorized
	Forbidden    = base.Forbidden
	NotFound     = base.NotFound
	Validation   = base.Validation
	Conflict     = base.Conflict
	Internal     = base.Internal
)
