package utility

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

var PlContext = &poolContext{
	limit:  16,
	holder: make(map[string]*WorkerPoolExecutor),
}

type poolContext struct {
	limit     int8
	workCount int8
	holder    map[string]*WorkerPoolExecutor
}

func newSingleton() *poolContext {
	if PlContext == nil {
		PlContext = &poolContext{
			limit:  16,
			holder: make(map[string]*WorkerPoolExecutor),
		}
		return PlContext
	}
	return PlContext
}

func (pc *poolContext) add(value *WorkerPoolExecutor) (key string, err error) {
	if pc.checkLimit() {
		return "", errors.New("work pool limit exceeded")
	}
	h := sha256.New()
	h.Write([]byte((value.id)))
	key = hex.EncodeToString(h.Sum(nil))
	if pc.holder[key] != nil {
		return "", errors.New("duplicated pool key: " + key)
	}
	pc.holder[key] = value
	pc.workCount++
	return key, nil
}

func (pc *poolContext) remove(key string) {
	delete(pc.holder, key)
}

func (pc *poolContext) checkLimit() bool {
	if pc.limit == pc.workCount {
		return true
	}
	return false
}

func (pc *poolContext) Get(key string) (*WorkerPoolExecutor, error) {
	var wp, ok = pc.holder[key]
	if !ok {
		return nil, errors.New("worker doesn't exist in the context with ID: " + key)
	}
	return wp, nil
}
