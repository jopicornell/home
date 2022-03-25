package models

import "time"

type Condition interface {
	IsLightOn(time time.Time) bool
	IsHeaterOn(time time.Time, temperature float64) bool
	IdealTemperature(time time.Time) float64
	GetMaxHours() float64
	GetMinHours() float64
}

type ConditionBase struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
