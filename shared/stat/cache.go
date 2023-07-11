
package stat

import "github.com/viant/gmetric/counter"

const (
	// Not deprecated since still used in other metrics
	NoSuchKey = "noKey"

	// Deprecated due to complications
	HasValue = "hasValue"
	CacheHit = "cacheHit"

	CacheCollision = "collision"
	CacheExpired   = "expired"
	LocalHasValue  = "localHasValue"
	LocalNoSuchKey = "localNoKey"