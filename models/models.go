package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base contains common columns for all tables
type Base struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;autoIncrement:false;type:uuid;not null"`
	CreatedAt time.Time `json:"createdAt" gorm:"type:timestamp without time zone;not null"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"type:timestamp without time zone;not null"`
}

// BeforeCreate will set Base struct before every insert
func (base *Base) BeforeCreate(tx *gorm.DB) error {
	// uuid.New() creates a new random UUID or panics.
	base.ID = uuid.New()

	// generate timestamps
	t := time.Now()
	base.CreatedAt, base.UpdatedAt = t, t

	return nil
}

// AfterUpdate will update the Base struct after every update
func (base *Base) AfterUpdate(tx *gorm.DB) error {
	// update timestamps
	base.UpdatedAt = time.Now()
	return nil
}
