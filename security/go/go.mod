module github.com/devpablocristo/core/security/go

go 1.26.1

require (
	github.com/devpablocristo/core/http/go v0.0.0
	github.com/google/uuid v1.6.0
)

replace github.com/devpablocristo/core/http/go => ../../http/go

replace github.com/devpablocristo/core/errors/go => ../../errors/go
