package model

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Tag struct {
	ID        uuid.UUID      `gorm:"column:id;primaryKey;default:uuidv7()"`
	Slug      string         `gorm:"column:slug"`
	Name      datatypes.JSON `gorm:"column:name;type:jsonb"`
	PackCount int            `gorm:"column:pack_count"`
}

func (Tag) TableName() string { return "tag" }

func (t *Tag) BeforeSave(*gorm.DB) error {
	t.Name = orEmptyJSON(t.Name)
	return nil
}

type PackTag struct {
	PackID uuid.UUID `gorm:"column:pack_id;primaryKey"`
	TagID  uuid.UUID `gorm:"column:tag_id;primaryKey"`
}

func (PackTag) TableName() string { return "pack_tag" }
