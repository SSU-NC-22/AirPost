package sql

import (
	"github.com/eunnseo/AirPost/application/domain/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type stationDroneRepo struct {
	db *gorm.DB
}

func NewStationDroneRepo() *stationDroneRepo {
	return &stationDroneRepo{
		db: dbConn,
	}
}

func (sdr *stationDroneRepo) FindsByStationID(stationid int) (sdl []model.StationDrone, err error) {
	return sdl, sdr.db.Where("station_id=?", stationid).Find(&sdl).Error
}

func (sdr *stationDroneRepo) Create(sd *model.StationDrone) error {
	return sdr.db.Omit(clause.Associations).Create(sd).Error
}

func (sdr *stationDroneRepo) Delete(sd *model.StationDrone) error {
	return sdr.db.Delete(sd).Error
}
