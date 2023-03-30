package resource

import (
	"sync"
)

var storage = struct {
	repos map[string]Repo
	mtx   sync.RWMutex
}{
	repos: make(map[string]Repo),
}

func GetRepo(name string) Repo {
	storage.mtx.RLock()
	defer storage.mtx.RUnlock()
	return storage.repos[name]
}

func registerRepo(resource string, repo Repo) {
	storage.mtx.Lock()
	defer storage.mtx.Unlock()
	storage.repos[resource] = repo
}
