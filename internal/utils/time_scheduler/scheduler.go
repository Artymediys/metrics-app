package time_scheduler

import (
	"fmt"
	"sort"
	"time"
)

func DurationUntilNextScheduledTime(scheduledTimes []string) (time.Duration, error) {
	now := time.Now()
	var durations []time.Duration

	for _, timeStr := range scheduledTimes {
		scheduledTime, err := time.Parse("15:04", timeStr)
		if err != nil {
			return 0, fmt.Errorf("utils.time_scheduler.DurationUntilNextScheduledTime -> time.Parse: %w", err)
		}

		nextScheduled := time.Date(
			now.Year(), now.Month(), now.Day(),
			scheduledTime.Hour(), scheduledTime.Minute(), 0, 0,
			now.Location(),
		)
		if nextScheduled.Before(now) {
			nextScheduled = nextScheduled.Add(24 * time.Hour)
		}

		durations = append(durations, nextScheduled.Sub(now))
	}

	sort.Slice(durations, func(i, j int) bool {
		return durations[i] < durations[j]
	})

	if len(durations) > 0 {
		return durations[0], nil
	}

	return 0, fmt.Errorf("utils.time_scheduler.DurationUntilNextScheduledTime: time not specified")
}
