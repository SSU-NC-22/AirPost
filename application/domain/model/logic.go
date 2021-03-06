package model

type Logic struct {
	ID		int		`json:"id" gorm:"primaryKey"`
	Name	string	`json:"name" gorm:"type:varchar(32);unique;not null"`
	Elems	string	`json:"elems" gorm:"type:text;not null"`
	NodeID	int		`json:"node_id" gorm:"not null"`
	Node	Node	`json:"node" gorm:"foreignKey:NodeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (Logic) TableName() string {
	return "logics"
}
