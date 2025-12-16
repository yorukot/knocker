package schedular

import (
	"math/rand"
	"time"
)

// calculateJitter returns a random duration between 0 and 30% of the interval
// If 30% of the interval is >= 20s, it caps the jitter at 20s
func calculateJitter(intervalSeconds int) time.Duration {
	// Calculate 30% of the interval
	thirtyPercent := float64(intervalSeconds) * 0.3

	// Cap at 20 seconds if 30% exceeds 20s
	maxJitterSeconds := thirtyPercent
	if thirtyPercent >= 5.0 {
		maxJitterSeconds = 5.0
	}

	// Generate random jitter between 0 and maxJitterSeconds
	jitterSeconds := rand.Float64() * maxJitterSeconds

	return time.Duration(jitterSeconds * float64(time.Second))
}
