package date

import (
	"fmt"
	"time"

	"github.com/yut-kt/goholiday"
)

func CountBusinessDay(startTime time.Time, endTime time.Time) (int, error) {
	if startTime.After(endTime) {
		return 0, fmt.Errorf("startTime is the time after endTime")
	}
	currentTime := startTime
	cnt := 0
	for currentTime.Before(endTime) {
		if goholiday.IsBusinessDay(currentTime) {
			cnt++
		}
		currentTime.Add(24 * time.Hour)
	}
	return cnt, nil
}
