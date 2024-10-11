package proxy

//go:generate go run github.com/dmarkham/enumer -type CacheStatus -trimprefix Cache -transform upper

type CacheStatus uint8

const (
	CacheMiss CacheStatus = iota
	CacheHit
	CacheBypass
)
