package models

import (
	"github.com/jopicornell/thermonick/utils"
	"time"
)

type SunsetSunriseCondition struct {
	MinimumHours     float64 `json:"minimumHours"`
	MaximumHours     float64 `json:"maximumHours"`
	TemperatureNight float64 `json:"temperatureNight"`
	TemperatureDay   float64 `json:"temperatureDay"`
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
}

func (c *SunsetSunriseCondition) IsHeaterOn(time time.Time, temperature float64) bool {
	sunrise, sunset := utils.AdjustedSunriseSunset(c.Latitude, c.Longitude, time, c.MinimumHours, c.MaximumHours)
	if time.After(sunrise) && time.Before(sunset) {
		return c.TemperatureDay-1 > temperature
	}
	return c.TemperatureNight-1 > temperature
}

func (c *SunsetSunriseCondition) IsHeaterOff(time time.Time, temperature float64) bool {
	sunrise, sunset := utils.AdjustedSunriseSunset(c.Latitude, c.Longitude, time, c.MinimumHours, c.MaximumHours)
	if time.After(sunrise) && time.Before(sunset) {
		return c.TemperatureDay < temperature
	}
	return c.TemperatureNight < temperature
}

func (c *SunsetSunriseCondition) IsLightOn(time time.Time) bool {
	sunrise, sunset := utils.AdjustedSunriseSunset(c.Latitude, c.Longitude, time, c.MinimumHours, c.MaximumHours)
	return time.After(sunrise) && time.Before(sunset)
}

func (c *SunsetSunriseCondition) IdealTemperature(time time.Time) float64 {
	sunrise, sunset := utils.AdjustedSunriseSunset(c.Latitude, c.Longitude, time, c.MinimumHours, c.MaximumHours)
	if time.After(sunrise) && time.Before(sunset) {
		return c.TemperatureDay
	}
	return c.TemperatureNight
}

func (c *SunsetSunriseCondition) GetMaxHours() float64 {
	return c.MaximumHours
}

func (c *SunsetSunriseCondition) GetMinHours() float64 {
	return c.MinimumHours
}
