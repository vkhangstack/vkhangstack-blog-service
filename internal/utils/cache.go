package utils

const (
	// CacheKeyPrefix is the prefix used for cache keys to avoid collisions.
	CacheKeyPrefix                   = "cache:"
	CacheKeyUserPrefix               = "user:"
	CacheKeyTemporarilyBlockedPrefix = "temp_blocked:"
	CacheKeyPostPrefix               = "post:"
	CacheKeyCategoryPrefix           = "category:"
)

type CacheTTL int64

const (
	// CacheTTLShort is the TTL for short-lived cache entries (1 minute).
	CacheTTLShort CacheTTL = 60

	// CacheTTLMid is the TTL for mid-lived cache entries (5 minutes).
	CacheTTLMid CacheTTL = 300

	// CacheTTLLong is the TTL for long-lived cache entries (1 hour).
	CacheTTLLong CacheTTL = 3600

	// CacheTTLForever is the TTL for cache entries that should never expire.
	CacheTTLForever CacheTTL = 0

	// CacheTTLOneMinute is the TTL for cache entries that should expire after one minute.
	CacheTTLOneMinute CacheTTL = 60

	// CacheTTLOneHour is the TTL for cache entries that should expire after one hour.
	CacheTTLOneHour CacheTTL = 3600

	// CacheTTLOneDay is the TTL for cache entries that should expire after one day.
	CacheTTLOneDay CacheTTL = 86400

	// CacheTTLOneWeek is the TTL for cache entries that should expire after one week.
	CacheTTLOneWeek CacheTTL = 604800
)
