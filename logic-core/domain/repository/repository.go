package repository

import (
	"github.com/eunnseo/AirPost/logic-core/domain/model"
)

type RegistRepo interface {
	FindNode(key int) (*model.Node, error)
	CreateNode(key int, n *model.Node) error
	DeleteNode(key int) error
	// FindSensor(key int) (*model.Sensor, error)
	// CreateSensor(key int, s *model.Sensor) error
	// DeleteSensor(key int) error
	AppendSinkAddr(sid int, s *string) error
}
