package cache

type ICache interface {
	Get(key interface{}) (value interface{}, ok bool)
	Put(key interface{}, value interface{})
	Delete(key interface{})
	Len() int
}
