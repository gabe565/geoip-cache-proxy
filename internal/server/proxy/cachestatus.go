package proxy

//go:generate enumer -type CacheStatus -trimprefix Cache -transform upper

type CacheStatus uint8

const (
	CacheMiss CacheStatus = iota
	CacheHit
	CacheBypass
)
