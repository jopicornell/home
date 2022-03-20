package temperatures

type Temperature struct {
	ID     string  `json:"id,omitempty" bson:"_id,omitempty"`
	Value  float64 `json:"value" bson:"value"`
	Sensor string  `json:"sensor" bson:"sensor"`
}
