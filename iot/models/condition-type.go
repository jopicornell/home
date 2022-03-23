package models

import (
	"bytes"
	"encoding/json"
)

type ConditionType string

const (
	IfLower  ConditionType = "IfLower"
	IfHigher ConditionType = "IfHigher"
	IfPast   ConditionType = "IfPast"
	IfFuture ConditionType = "IfFuture"
)

var TemperatureConditionMap = map[ConditionType]string{
	IfLower:  string(IfLower),
	IfHigher: string(IfHigher),
}

// MarshalJSON marshals the enum as a quoted json string
func (s ConditionType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(TemperatureConditionMap[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *ConditionType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = ConditionType(j)
	return nil
}
