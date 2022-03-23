package models

import (
	"github.com/jopicornell/thermonick/utils"
	"time"
)

type SunsetSunriseCondition struct {
	Condition
	MinimumHours     float32 `json:"minimumHours"`
	MaximumHours     float32 `json:"maximumHours"`
	TemperatureNight float32 `json:"temperatureNight"`
	TemperatureDay   float32 `json:"temperatureDay"`
}

func (c *SunsetSunriseCondition) IsHeaterOn(time time.Time, temperature float32) bool {
	sunrise, sunset := utils.SunsetSunrise(39.57, 2.65, time)
	if sunrise.After(time) || sunset.Before(time) {
		return c.TemperatureNight > temperature
	}
	return c.TemperatureDay > temperature
}

func (c *SunsetSunriseCondition) IsLightOn(time time.Time) bool {
	sunrise, sunset := utils.SunsetSunrise(39.57, 2.65, time)
	if sunrise.After(time) || sunset.Before(time) {
		return false
	}
	return true
}
