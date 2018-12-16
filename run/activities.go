package run

import (
	"fmt"
	"time"

	"github.com/globalsign/mgo/bson"
)

func (v *VirtualRun) AthleteSummary(id uint32) (*Stats, error) {
	stats := &Stats{}
	activities := []Activity{}

	now := time.Now()
	currentYear, currentMonth, currentDay := now.Date()
	currentLocation := now.Location()

	firstOfYear := time.Date(currentYear, 1, 1, 0, 0, 0, 0, currentLocation)
	firstOfYearStamp := firstOfYear.Format(time.RFC3339)

	// stats only activities in this year
	e := v.db.List("activities", bson.M{
		"athlete": bson.M{
			"id": id,
		},
		"startdate": bson.M{
			"$gte": firstOfYearStamp,
		},
	}, &activities)

	if e != nil {
		return nil, e
	}

	for _, activity := range activities {
		dateTaken, _ := time.Parse(time.RFC3339, activity.StartDate)
		offset, _ := time.ParseDuration(fmt.Sprintf("%fs", activity.TimeZoneOffset))
		dateTaken = dateTaken.Add(offset)

		stats.ThisYearRunTotals.Count++
		stats.ThisYearRunTotals.Distance += activity.Distance
		stats.ThisYearRunTotals.ElapsedTime += activity.ElapsedTime
		stats.ThisYearRunTotals.ElevationGain += activity.ElevationGain
		stats.ThisYearRunTotals.MovingTime += activity.MovingTime

		if currentMonth == dateTaken.Month() {
			stats.ThisMonthRunTotals.Count++
			stats.ThisMonthRunTotals.Distance += activity.Distance
			stats.ThisMonthRunTotals.ElapsedTime += activity.ElapsedTime
			stats.ThisMonthRunTotals.ElevationGain += activity.ElevationGain
			stats.ThisMonthRunTotals.MovingTime += activity.MovingTime

			if currentDay == dateTaken.Day() {
				stats.RecentRun.Distance = activity.Distance
				stats.RecentRun.Title = activity.Title
				stats.RecentRun.ElapsedTime = activity.ElapsedTime
				stats.RecentRun.MovingTime = activity.MovingTime
				stats.RecentRun.StartDate = activity.StartDate
				stats.RecentRun.TimeZoneOffset = activity.TimeZoneOffset
			}
		}
	}

	return stats, nil
}
