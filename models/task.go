package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// RecurringPattern defines how a task should recur
type RecurringPattern struct {
	Type     string `json:"type"`      // once, daily, weekly, monthly, custom
	Interval int    `json:"interval"`  // Repeat every X days/weeks/months
	Weekdays []int  `json:"weekdays"`  // 0-6 for Sunday-Saturday
	EndDate  string `json:"end_date"`  // When to stop recurring
	CronExpr string `json:"cron_expr"` // Custom cron expression
}

// TaskProgress represents the progress of a task
type TaskProgress struct {
	Percentage int       `json:"percentage"`
	Status     string    `json:"status"`
	UpdatedAt  time.Time `json:"updated_at"`
	Message    string    `json:"message"`
}

type Task struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Name         string         `json:"name" gorm:"not null"`
	Description  string         `json:"description"`
	ScheduleTime time.Time      `json:"scheduleTime" gorm:"not null;index"`
	Priority     string         `json:"priority" gorm:"type:varchar(10);check:priority in ('low', 'medium', 'high')"`
	Status       string         `json:"status" gorm:"type:varchar(20);default:'pending';check:status in ('pending', 'running', 'completed', 'failed')"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
	RetryCount   int            `json:"retryCount" gorm:"default:0"`
	LastError    string         `json:"lastError"`
	Tags         []string       `json:"tags" gorm:"type:text[]"`
	Metadata     string         `json:"metadata" gorm:"type:jsonb"`

	// New fields
	ParentTaskID    *uint           `json:"parentTaskId" gorm:"index"` // For task dependencies
	DependentTasks  []Task          `json:"dependentTasks" gorm:"foreignKey:ParentTaskID"`
	IsRecurring     bool            `json:"isRecurring" gorm:"default:false"`
	RecurringConfig json.RawMessage `json:"recurringConfig" gorm:"type:jsonb"` // Stores RecurringPattern
	Progress        TaskProgress    `json:"progress" gorm:"embedded"`
	Assignee        string          `json:"assignee"` // User assigned to the task
	DueDate         *time.Time      `json:"dueDate"`
	EstimatedTime   int             `json:"estimatedTime"`                  // In minutes
	ActualTime      int             `json:"actualTime"`                     // In minutes
	Labels          []string        `json:"labels" gorm:"type:text[]"`      // For better organization
	Priority_Score  float64         `json:"priorityScore" gorm:"default:0"` // Calculated priority score
	Comments        []TaskComment   `json:"comments" gorm:"foreignKey:TaskID"`
	Attachments     []Attachment    `json:"attachments" gorm:"foreignKey:TaskID"`
}

type TaskComment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	TaskID    uint      `json:"taskId"`
	UserID    string    `json:"userId"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Attachment struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	TaskID      uint      `json:"taskId"`
	FileName    string    `json:"fileName"`
	FileType    string    `json:"fileType"`
	FileSize    int64     `json:"fileSize"`
	StoragePath string    `json:"storagePath"`
	UploadedBy  string    `json:"uploadedBy"`
	CreatedAt   time.Time `json:"createdAt"`
}
