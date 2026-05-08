package worker

import (
	"math/rand/v2"
	"time"
)

const jitterFraction = 0.2

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

	jitter := float64(delay) * jitterFraction * (rand.Float64()*2 - 1)
	result := time.Duration(float64(delay) + jitter)

	if result > max {
		result = max
	}
	if result <= 0 {
		result = base
	}
	return result
}
