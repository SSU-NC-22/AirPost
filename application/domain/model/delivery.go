package model

import "time"

type Delivery struct {
	ID			 int	   `json:"id" gorm:"primaryKey"`
	OrderNum	 string	   `json:"order_num" gorm:"type:varchar(32);not null"`
	DroneID		 int	   `json:"drone_id" gorm:"not null"`
	Drone		 Node	   `json:"drone_node" gorm:"foreignKey:DroneID"`
	Email		 string	   `json:"email" gorm:"type:varchar(32);not null"`
	SrcName		 string	   `json:"src_name" gorm:"type:varchar(32);not null"`
	SrcPhone	 string	   `json:"src_phone" gorm:"type:varchar(32);not null"`
	SrcStationID int	   `json:"src_station_id" gorm:"not null"`
	DestName	 string	   `json:"dest_name" gorm:"type:varchar(32);not null"`
	DestPhone	 string	   `json:"dest_phone" gorm:"type:varchar(32);not null"`
	DestTagID    int	   `json:"dest_tag_id" gorm:"not null"`
	CreatedAt	 time.Time `json:"created_at" gorm:"not null"`
}

func (Delivery) TableName() string {
	return "deliveries"
}

type StationDrone struct{
	ID		  int  `json:"id" gorm:"primaryKey"`
	StationID int  `json:"station_id" gorm:"not null"`
	Station	  Node `json:"station" gorm:"foreignKey:StationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	DroneID   int  `json:"drone_id" gorm:"not null"`
	Drone	  Node `json:"drone" gorm:"foreignKey:DroneID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Usable    bool `json:"usable" gorm:"not null"` // true : 사용가능, false : 불가능
}

func (StationDrone) TableName() string {
	return "station_drone"
}

type Path struct {
	ID		  int	  `json:"id" gorm:"primaryKey"`
	StationID int	  `json:"station_id" gorm:"not null"`
	TagID 	  int	  `json:"tag_id" gorm:"not null"`
	Path   	  string  `json:"path" gorm:"type:text;not null"`
	Distance  float64 `json:"distance" gorm:"not null"`
}

func (Path) TableName() string {
	return "paths"
}

type Tracking struct {
	DroneNid int     `json:"droneNid"`
	SrcLat	 float64 `json:"srcLat"`
	SrcLng	 float64 `json:"srcLng"`
	DestLat	 float64 `json:"destLat"`
	DestLng	 float64 `json:"destLng"`
	DroneLat float64 `json:"droneLat"`
	DroneLng float64 `json:"droneLng"`
}
