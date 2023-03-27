package resource

// access control
type AC interface {
	WhoAmI() uint
	Allow(action string, resource string, record interface{}) bool
}
