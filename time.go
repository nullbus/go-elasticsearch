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

func (d Duration) MarshalJSON() ([]byte, error) {
	str := `"` + d.String() + `"`
	return []byte(str), nil
}

func (d Duration) String() string {
	var minusSign string
	if d < 0 {
		d, minusSign = -d, "-"
	}

	if d >= Year {
		return fmt.Sprintf("%s%fy", minusSign, float64(d)/float64(Year))
	}
	if d >= Month {
		return fmt.Sprintf("%s%fM", minusSign, float64(d)/float64(Month))
	}
	if d >= Week {
		return fmt.Sprintf("%s%fw", minusSign, float64(d)/float64(Week))
	}
	if d >= Day {
		return fmt.Sprintf("%s%fd", minusSign, float64(d)/float64(Day))
	}
	if d >= Hour {
		return fmt.Sprintf("%s%fh", minusSign, float64(d)/float64(Hour))
	}
	if d >= Minute {
		return fmt.Sprintf("%s%fm", minusSign, float64(d)/float64(Minute))
	}

	return fmt.Sprintf("%s%fs", minusSign, float64(d)/float64(Second))
}
