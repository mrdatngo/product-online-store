package model

import (
	"gorm.io/gorm"
	"time"
)

type EntityBranch struct {
	ID        int64     `json:"id" gorm:"type:uint;primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (entity *EntityBranch) BeforeCreate(db *gorm.DB) error {
	entity.CreatedAt = time.Now().Local()
	return nil
}

func (entity *EntityBranch) BeforeUpdate(db *gorm.DB) error {
	entity.UpdatedAt = time.Now().Local()
	return nil
}

// TableName return table name type string to define sql table
func (EntityBranch) TableName() string {
	return "branch"
}
