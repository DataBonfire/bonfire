package resource

import (
	"sync"
)

var (
	storage    = make(map[string]Repo)
	storageMtx sync.Mutex
)

func registerRepo(resource string, repo Repo) {
	storageMtx.Lock()
	defer storageMtx.Unlock()
	storage[resource] = repo
}
