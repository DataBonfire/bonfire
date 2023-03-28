package resource

import "errors"

// access control
type AC interface {
	WhoAmI() uint
	Allow(action string, resource string, record interface{}) bool
}

var (
	ErrPermissionDenied = errors.New("permission denied")
)
