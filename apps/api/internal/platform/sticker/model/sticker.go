package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Sticker struct {
	ID            uuid.UUID      `gorm:"column:id;primaryKey;default:uuidv7()"`
	PackID        uuid.UUID      `gorm:"column:pack_id"`
	Position      int            `gorm:"column:position"`
	ImageHash     string         `gorm:"column:image_hash"`
	Width         int            `gorm:"column:width"`
	Height        int            `gorm:"column:height"`
	Game          datatypes.JSON `gorm:"column:game;type:jsonb"`
	CharacterName datatypes.JSON `gorm:"column:character_name;type:jsonb"`
	VndbID        *int           `gorm:"column:vndb_id"`
	Note          string         `gorm:"column:note"`
	CreatedAt     time.Time      `gorm:"column:created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at"`
}

func (Sticker) TableName() string { return "sticker" }

func (s *Sticker) BeforeSave(*gorm.DB) error {
	s.Game = orEmptyJSON(s.Game)
	s.CharacterName = orEmptyJSON(s.CharacterName)
	return nil
}
