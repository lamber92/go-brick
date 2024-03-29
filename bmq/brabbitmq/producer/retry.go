package producer

import "time"

const _defaultMaxRetryTimes = 2

// defaultRetryTimeInterval
func defaultRetryTimeInterval(times uint) (sec time.Duration) {
	switch times {
	case 1:
		sec = time.Second
	case 2:
		sec = 2 * time.Second
	case 3:
		sec = 4 * time.Second
	case 4:
		sec = 16 * time.Second
	default:
		sec = time.Minute
	}
	return sec
}
