package models

import (
	"github.com/jopicornell/thermonick/utils"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var dayTimes = []time.Time{
	time.Date(2019, time.December, 21, 9, 6, 0, 0, time.UTC),
	time.Date(2019, time.December, 21, 16, 55, 0, 0, time.UTC),
	time.Date(2019, time.June, 24, 7, 1, 0, 0, time.UTC),
	time.Date(2019, time.June, 21, 18, 55, 0, 0, time.UTC),
}

func TestSunsetSunriseDay(t *testing.T) {
	condition := &SunsetSunriseCondition{
		MinimumHours:     8,
		MaximumHours:     12,
		TemperatureNight: 20,
		TemperatureDay:   30,
		Latitude:         39.57,
		Longitude:        2.65,
	}
	for i, testTime := range dayTimes {
		t.Run(testTime.String(), func(t *testing.T) {
			testTime := dayTimes[i]
			t.Parallel()
			sunrise, sunset := utils.AdjustedSunriseSunset(condition.Latitude, condition.Longitude, testTime, condition.MinimumHours, condition.MaximumHours)
			log.Printf("Testing %s %s %s", testTime.String(), sunrise.String(), sunset.String())
			assert.Truef(t, condition.IsHeaterOn(testTime, 20), "Heater should be on")
			assert.Falsef(t, condition.IsHeaterOn(testTime, 31), "Heater should be off")
			assert.Falsef(t, condition.IsHeaterOn(testTime, 40), "Heater should be off")
			assert.Falsef(t, condition.IsHeaterOn(testTime, condition.TemperatureDay), "Heater should be off")
			assert.Falsef(t, condition.IsHeaterOn(testTime, 29.1), "Heater should be off")
			assert.Falsef(t, condition.IsHeaterOn(testTime, 29), "Heater should be off")
			assert.Truef(t, condition.IsHeaterOn(testTime, 28.9), "Heater should be on")

			assert.Falsef(t, condition.IsHeaterOff(testTime, 20), "HeaterOff checker should be false when temperature (%d) is below or equal to %f", 20, condition.TemperatureDay)
			assert.Falsef(t, condition.IsHeaterOff(testTime, condition.TemperatureDay), "HeaterOff checker should be false when temperature (%d) is below or equal to %f", 30, condition.TemperatureDay)
			assert.True(t, condition.IsHeaterOff(testTime, 30.1), "HeaterOff checker should be true when temperature (%f) is above %f", 30.1, condition.TemperatureDay)
			assert.True(t, condition.IsHeaterOff(testTime, 40), "HeaterOff checker should be true when temperature (%d) is above %f", 40, condition.TemperatureDay)

			assert.True(t, condition.IsLightOn(testTime), "Expected light to be on")
			assert.Equalf(t, condition.IdealTemperature(testTime), condition.TemperatureDay, "Expected ideal temperature to be %d", condition.TemperatureDay)
		})
	}
}

var nightTimes = []time.Time{
	time.Date(2019, time.December, 21, 8, 0, 0, 0, time.UTC),
	time.Date(2019, time.December, 21, 18, 0, 0, 0, time.UTC),
	time.Date(2019, time.December, 21, 3, 0, 0, 0, time.UTC),
	time.Date(2019, time.December, 21, 21, 0, 0, 0, time.UTC),
	time.Date(2019, time.June, 24, 6, 59, 0, 0, time.UTC),
	time.Date(2019, time.June, 21, 19, 1, 0, 0, time.UTC),
	time.Date(2019, time.June, 24, 3, 0, 0, 0, time.UTC),
	time.Date(2019, time.June, 21, 22, 0, 0, 0, time.UTC),
}

func TestSunsetSunriseNight(t *testing.T) {
	for i, testTime := range nightTimes {
		t.Run(testTime.String(), func(t *testing.T) {
			testTime := nightTimes[i]
			t.Parallel()
			condition := &SunsetSunriseCondition{
				MinimumHours:     8,
				MaximumHours:     12,
				TemperatureNight: 21,
				TemperatureDay:   30,
				Latitude:         39.57,
				Longitude:        2.65,
			}
			assert.Truef(t, condition.IsHeaterOn(testTime, 19.9), "Heater should be on")
			assert.Truef(t, condition.IsHeaterOn(testTime, 10), "Heater should be on")
			assert.Falsef(t, condition.IsHeaterOn(testTime, 31), "Heater should be off")
			assert.Falsef(t, condition.IsHeaterOn(testTime, 40), "Heater should be off")
			assert.Falsef(t, condition.IsHeaterOn(testTime, condition.TemperatureNight), "Heater should be off")
			assert.Falsef(t, condition.IsHeaterOn(testTime, 29.1), "Heater should be off")
			assert.Falsef(t, condition.IsHeaterOn(testTime, 29), "Heater should be off")

			assert.Falsef(t, condition.IsHeaterOff(testTime, 20), "HeaterOff checker should be false when temperature (%d) is below or equal to %f", 20, condition.TemperatureNight)
			assert.Falsef(t, condition.IsHeaterOff(testTime, condition.TemperatureNight), "HeaterOff checker should be false when temperature (%d) is below or equal to %f", condition.TemperatureNight, condition.TemperatureNight)
			assert.True(t, condition.IsHeaterOff(testTime, 21.1), "HeaterOff checker should be true when temperature (%f) is above %f", 30.1, condition.TemperatureNight)
			assert.True(t, condition.IsHeaterOff(testTime, 30), "HeaterOff checker should be true when temperature (%d) is above %f", 40, condition.TemperatureNight)

			assert.False(t, condition.IsLightOn(testTime), "Expected light to be on")
			assert.Equalf(t, condition.IdealTemperature(testTime), condition.TemperatureNight, "Expected ideal temperature to be %d", condition.TemperatureNight)
		})
	}
}
