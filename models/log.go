package models

import (
	"gorm.io/datatypes"
)

type Log struct {
	ID       int `gorm:"primaryKey"`
	DeviceID string
	Path     string
	Data     datatypes.JSON
}
