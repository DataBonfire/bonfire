package ac

type Accessor interface {
	GetID() uint
	GetGroups() []uint
	GetSubordinates() []uint
}
