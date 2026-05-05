package worker

import (
	"math/rand/v2"
	"time"
)

func NextBackoff(attempt int, base, max time.Duration) time.Duration {
	if attempt < 1 {
		attempt = 1
	}

	delay := base
	for i := 1; i < attempt; i++ {
		delay *= 2
		if delay >= max || delay <= 0 {
			delay = max
			break
		}
	}
	if delay > max || delay <= 0 {
		delay = max
	}
	if delay <= 0 {
		return 0
	}

	return time.Duration(rand.Int64N(int64(delay)))
}
