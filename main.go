package main

import (
	"fmt"
	"math"
	"time"
)

func main() {
	startTime := time.Now()

	var durations []time.Duration

	durations = append(durations, time.Since(startTime)+47*time.Second+1243*time.Millisecond)
	durations = append(durations, time.Since(startTime)+18*time.Minute+36*time.Second+1243*time.Millisecond)
	durations = append(durations, time.Since(startTime)+2*time.Hour+15*time.Minute+36*time.Second+1243*time.Millisecond)

	for _, duration := range durations {
		fmt.Println(duration)
		fmt.Println(fmtDuration(duration))
	}
}

func fmtDuration(d time.Duration) string {
	numSeconds := int(d.Seconds())
	s := numSeconds % 60
	m := int(math.Floor(float64(numSeconds)/60)) % 60
	h := math.Floor(float64(numSeconds) / 3600)

	if h > 0 {
		return fmt.Sprintf("%vh:%vm:%vs", h, m, s)
	} else if m > 0 {
		return fmt.Sprintf("%vm:%vs", m, s)
	} else {
		return fmt.Sprintf("%vs", s)
	}
}
