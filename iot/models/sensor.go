package models

type Sensor struct {
	Id         string
	Name       string
	Value      float32
	Conditions []Condition
}
