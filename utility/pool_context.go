package utility

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

var plContext *poolContext

type poolContext struct {
	limit     int8
	workCount int8
	holder    map[string]*WorkerPoolExecutor
}

func newSingleton() *poolContext {
	if plContext == nil {
		plContext = &poolContext{
			limit:  16,
			holder: make(map[string]*WorkerPoolExecutor),
		}
		return plContext
	}
	return plContext
}

func (pc *poolContext) Add(value *WorkerPoolExecutor) (key string, err error) {
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

func (pc *poolContext) Remove(key string) {
	delete(pc.holder, key)
}

func (pc *poolContext) checkLimit() bool {
	if pc.limit == pc.workCount {
		return true
	}
	return false
}
