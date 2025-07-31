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

func (pc *poolContext) generateWorkerKey(value *WorkerPoolExecutor) string {
	h := sha256.New()
	h.Write([]byte((value.name)))
	key := hex.EncodeToString(h.Sum(nil))
	return key
}

func (pc *poolContext) add(value *WorkerPoolExecutor) error {
	if pc.checkLimit() {
		return errors.New("work pool limit exceeded")
	}

	if pc.holder[value.ID] != nil {
		return errors.New("duplicated pool key: " + value.ID)
	}
	pc.holder[value.ID] = value
	pc.workCount++
	return nil
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
