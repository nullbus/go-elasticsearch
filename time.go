package es

import (
	"fmt"
	"time"
)

type Duration time.Duration

const (
	Second = Duration(time.Second)
	Minute = Duration(time.Minute)
	Hour   = Duration(time.Hour)
	Day    = Duration(24 * time.Hour)
	Week   = Duration(7 * Day)
	Month  = Duration(30 * Day)
	Year   = Duration(365 * Day)
)

func (d Duration) String() string {
	if d >= Year {
		return fmt.Sprintf("%fy", float64(d)/float64(Year))
	}
	if d >= Month {
		return fmt.Sprintf("%fM", float64(d)/float64(Month))
	}
	if d >= Week {
		return fmt.Sprintf("%fw", float64(d)/float64(Week))
	}
	if d >= Day {
		return fmt.Sprintf("%fd", float64(d)/float64(Day))
	}
	if d >= Hour {
		return fmt.Sprintf("%fh", float64(d)/float64(Hour))
	}
	if d >= Minute {
		return fmt.Sprintf("%fm", float64(d)/float64(Minute))
	}

	return fmt.Sprintf("%fs", float64(d)/float64(Second))
}
