package utils

import (
	"github.com/jopicornell/thermonick/log"
	"github.com/kelvins/sunrisesunset"
	"time"
)

var backupSunset = time.Date(2009, time.November, 10, 19, 0, 0, 0, time.UTC)
var backupSunrise = time.Date(2009, time.November, 10, 8, 30, 0, 0, time.UTC)

func SunsetSunrise(latitude float64, longitude float64, date time.Time) (sunrise time.Time, sunset time.Time) {
	_, offset := date.Zone()
	log.Logger.Debugf("Offset: %d", offset)
	sunrise, sunset, err := sunrisesunset.GetSunriseSunset(latitude, longitude, float64(offset), date)
	if err != nil {
		log.Logger.Error(err)
		sunrise = backupSunrise
		sunset = backupSunset
	}
	return MergeDateTime(date, sunrise), MergeDateTime(date, sunset)
}

func MergeDateTime(date time.Time, onlyTime time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), onlyTime.Hour(), onlyTime.Minute(), onlyTime.Second(), onlyTime.Nanosecond(), date.Location())
}
