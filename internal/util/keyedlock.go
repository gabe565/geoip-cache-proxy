package util

import "sync"

type KeyedLock struct {
	mu      sync.Mutex
	keyedMu map[string]*sync.Mutex
}

func NewKeyedLock() *KeyedLock {
	return &KeyedLock{
		keyedMu: make(map[string]*sync.Mutex),
	}
}

func (l *KeyedLock) getLockBy(key string) *sync.Mutex {
	l.mu.Lock()
	defer l.mu.Unlock()

	if mu, ok := l.keyedMu[key]; ok {
		return mu
	}

	mu := &sync.Mutex{}
	l.keyedMu[key] = mu
	return mu
}

func (l *KeyedLock) Lock(key string) { l.getLockBy(key).Lock() }

func (l *KeyedLock) Unlock(key string) { l.getLockBy(key).Unlock() }
