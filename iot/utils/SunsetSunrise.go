package utils

import (
	"github.com/jopicornell/thermonick/log"
	"github.com/kelvins/sunrisesunset"
	"time"
)

var backupSunset = time.Date(2009, time.November, 10, 19, 0, 0, 0, time.UTC)
var backupSunrise = time.Date(2009, time.November, 10, 8, 30, 0, 0, time.UTC)
var MidTime = time.Date(2009, time.November, 10, 13, 0, 0, 0, time.UTC)

func SunsetSunrise(latitude float64, longitude float64, date time.Time) (sunrise time.Time, sunset time.Time) {
	_, offset := date.Zone()
	sunrise, sunset, err := sunrisesunset.GetSunriseSunset(latitude, longitude, float64(offset/int(time.Hour.Seconds())), date)
	if err != nil {
		log.Logger.Errorf("Error: %+v", err)
		sunrise = backupSunrise
		sunset = backupSunset
	}
	return MergeDateTime(date, sunrise), MergeDateTime(date, sunset)
}

func AdjustedSunriseSunset(latitude float64, longitude float64, date time.Time, maxHours float64, minHours float64) (sunrise time.Time, sunset time.Time) {
	sunrise, sunset = SunsetSunrise(latitude, longitude, date)
	sunrise = sunrise.Add(1 * time.Hour)
	maxDuration := time.Duration(maxHours * float64(time.Hour.Nanoseconds()))
	minDuration := time.Duration(minHours * float64(time.Hour.Nanoseconds()))
	dayDuration := MinMaxDuration(GetDayDuration(sunrise, sunset), minDuration, maxDuration)
	halfDayDuration := dayDuration / 2
	sunrise = MidTime.Add(-halfDayDuration)
	sunset = MidTime.Add(halfDayDuration)
	return MergeDateTime(date, sunrise), MergeDateTime(date, sunset)
}

func GetDayDuration(date1, date2 time.Time) time.Duration {
	if date1.After(date2) {
		return date1.Sub(date2)
	}
	return date2.Sub(date1)
}

func MinMaxDuration(duration time.Duration, minDuration time.Duration, maxDuration time.Duration) time.Duration {
	if duration > maxDuration {
		log.Logger.Infof("Duration %+v is greater than max duration %+v", duration, maxDuration)
		return maxDuration
	}
	if duration < minDuration {
		log.Logger.Infof("Duration %+v is less than min duration %+v", duration, minDuration)
		return minDuration
	}
	return duration
}

func MergeDateTime(date time.Time, onlyTime time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), onlyTime.Hour(), onlyTime.Minute(), onlyTime.Second(), onlyTime.Nanosecond(), date.Location())
}
