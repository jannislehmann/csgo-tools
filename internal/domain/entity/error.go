package entity

import "errors"

var ErrNotFound = errors.New("infrastructure: entity not found")

var ErrInvalidEntity = errors.New("infrastructure: invalid entity")

var ErrCannotBeDeleted = errors.New("infrastructure: cannot be deleted")

var ErrUnknownInfrastructureError = errors.New("infrastructure: unknown error during database operation")
