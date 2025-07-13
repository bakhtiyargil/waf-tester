package utility

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

var plContext *poolContext

type poolContext struct {
	workCount int8
	holder    map[string]*WorkerPoolExecutor
}

func newSingleton() *poolContext {
	if plContext == nil {
		return &poolContext{
			holder: make(map[string]*WorkerPoolExecutor),
		}
	}
	return plContext
}

func (pc *poolContext) Add(value *WorkerPoolExecutor) (key string, err error) {
	if pc.holder[key] != nil {
		return "", errors.New("duplicated pool key: " + key)
	}
	h := sha256.New()
	h.Write([]byte((value.id)))
	key = hex.EncodeToString(h.Sum(nil))
	pc.holder[key] = value
	return key, nil
}

func (pc *poolContext) Remove(key string) {
	delete(pc.holder, key)
}
