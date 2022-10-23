package model

import (
	"gorm.io/gorm"
	"time"
)

type EntityDetailLog struct {
	DataType  string    `json:"event_type" gorm:"type:varchar(255);not null"`
	Param     string    `json:"param" gorm:"type:varchar(255);"`
	Value     string    `json:"value" gorm:"type:varchar(255);"`
	CreatedAt time.Time `json:"-"`
}

func (entity *EntityDetailLog) BeforeCreate(db *gorm.DB) error {
	entity.CreatedAt = time.Now().Local()
	return nil
}

func (EntityDetailLog) TableName() string {
	return "user_log_detail"
}

type EntityUserLog struct {
	UserID     int64              `json:"user_id" gorm:"type:uint"`
	EventType  string             `json:"event_type" gorm:"type:varchar(255);not null"`
	CreatedAt  time.Time          `json:"-"`
	DetailLogs []*EntityDetailLog `json:"detail_logs"`
}

func (entity *EntityUserLog) BeforeCreate(db *gorm.DB) error {
	entity.CreatedAt = time.Now().Local()
	return nil
}

func (EntityUserLog) TableName() string {
	return "user_log"
}
