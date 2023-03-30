package resource

import "errors"

// access control
type AC interface {
	Whoami() uint
	Allow(action string, resource string, record interface{}) bool
	GetFilters(action string, resource string) []Filter
}

var (
	ErrPermissionDenied = errors.New("permission denied")
)
