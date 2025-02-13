package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Name         string         `json:"name" gorm:"not null"`
	Description  string         `json:"description"`
	ScheduleTime time.Time      `json:"scheduleTime" gorm:"not null"`
	Priority     string         `json:"priority" gorm:"type:varchar(10);check:priority in ('low', 'medium', 'high')"`
	Status       string         `json:"status" gorm:"type:varchar(20);default:'pending';check:status in ('pending', 'running', 'completed', 'failed')"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
	RetryCount   int            `json:"retryCount" gorm:"default:0"`
	LastError    string         `json:"lastError"`
	Tags         []string       `json:"tags" gorm:"type:text[]"`
	Metadata     string         `json:"metadata" gorm:"type:jsonb"`
} 