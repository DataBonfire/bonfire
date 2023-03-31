package ac

import (
	"errors"

	"github.com/databonfire/bonfire/filter"
)

type AccessController interface {
	Allow(accessor interface{}, action string, resource string, record interface{}) bool
	Filters(accessor interface{}, action string, resource string) []filter.Filter
}

var (
	ErrUnknownAccessController = errors.New("unknown access controller")
	ErrPermissionDenied        = errors.New("permission denied")
)
