package run

import (
	"fmt"
	"time"

	"github.com/globalsign/mgo/bson"
)

func (v *VirtualRun) AthleteActivities(id uint32, period []string) ([]Activity, error) {
	activities := []Activity{}

	e := v.db.List("activities", bson.M{
		"athlete": bson.M{
			"id": id,
		},
		"startdate": bson.M{
			"$gte": period[0],
			"$lte": period[1],
		},
	}, []string{"-startdate"}, &activities)

	if e != nil {
		return nil, e
	}

	return activities, nil
}

func (v *VirtualRun) AthleteSummary(id uint32) (*Stats, error) {
	stats := &Stats{}
	activities := []Activity{}

	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
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
	}, []string{"-startdate"}, &activities)

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

		}
	}

	if len(activities) > 0 {
		stats.RecentRun.Distance = activities[0].Distance
		stats.RecentRun.Title = activities[0].Title
		stats.RecentRun.ElapsedTime = activities[0].ElapsedTime
		stats.RecentRun.MovingTime = activities[0].MovingTime
		stats.RecentRun.StartDate = activities[0].StartDate
		stats.RecentRun.TimeZoneOffset = activities[0].TimeZoneOffset
	}

	return stats, nil
}
