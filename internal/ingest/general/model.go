package general

import "time"

// SensorData defines the structure for "general" data points.
type SensorData struct {
	Time       time.Time `json:"time"`
	ShipID     string    `json:"ship_id"   validate:"required,min=1,max=100"`
	CargoID    string    `json:"cargo_id"  validate:"required,min=1,max=100"`
	Value      *float64   `json:"value"     validate:"required"`
}
