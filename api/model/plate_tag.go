package model

import "github.com/google/uuid"

type PlateTag struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey"      json:"id"`
	PlateID uuid.UUID `gorm:"type:uuid;not null;index"  json:"plate_id"`
	Tag     string    `gorm:"size:100;not null"         json:"tag"`

	Plate *Plate `gorm:"foreignKey:PlateID" json:"-"`
}

func (PlateTag) TableName() string {
	return "plate_tag"
}
