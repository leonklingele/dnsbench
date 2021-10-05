package dnsbench

import (
	"time"
)

type Measurement struct {
	Duration time.Duration
}

type Measurements []Measurement

func (ms Measurements) AverageDuration() time.Duration {
	var sum time.Duration
	for _, m := range ms {
		sum += m.Duration
	}

	return time.Duration(sum.Nanoseconds() / int64(len(ms)))
}
