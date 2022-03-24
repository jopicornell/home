package models

import (
	"github.com/jopicornell/thermonick/utils"
	"time"
)

type SunsetSunriseCondition struct {
	Condition
	MinimumHours     float64 `json:"minimumHours"`
	MaximumHours     float64 `json:"maximumHours"`
	TemperatureNight float64 `json:"temperatureNight"`
	TemperatureDay   float64 `json:"temperatureDay"`
}

func (c *SunsetSunriseCondition) IsHeaterOn(time time.Time, temperature float64) bool {
	sunrise, sunset := utils.SunsetSunrise(39.57, 2.65, time)
	if time.After(sunrise) && time.Before(sunset) {
		return c.TemperatureDay > temperature
	}
	return c.TemperatureNight > temperature
}

func (c *SunsetSunriseCondition) IsLightOn(time time.Time) bool {
	sunrise, sunset := utils.SunsetSunrise(39.57, 2.65, time)
	return time.After(sunrise) && time.Before(sunset)
}

func (c *SunsetSunriseCondition) IdealTemperature(time time.Time) float64 {
	sunrise, sunset := utils.SunsetSunrise(39.57, 2.65, time)
	if time.After(sunrise) && time.Before(sunset) {
		return c.TemperatureDay
	}
	return c.TemperatureNight
}
