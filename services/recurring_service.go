package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/task-schedulart/models"
	"gorm.io/gorm"
)

type RecurringTaskService struct {
	db   *gorm.DB
	cron *cron.Cron
}

func NewRecurringTaskService(db *gorm.DB) *RecurringTaskService {
	return &RecurringTaskService{
		db:   db,
		cron: cron.New(cron.WithSeconds()),
	}
}

// StartScheduler starts the recurring task scheduler
func (s *RecurringTaskService) StartScheduler() error {
	// Start the cron scheduler
	s.cron.Start()

	// Load all recurring tasks from database
	var tasks []models.Task
	if err := s.db.Where("is_recurring = ?", true).Find(&tasks).Error; err != nil {
		return fmt.Errorf("failed to load recurring tasks: %v", err)
	}

	// Schedule each task
	for _, task := range tasks {
		if err := s.scheduleTask(task); err != nil {
			return fmt.Errorf("failed to schedule task %d: %v", task.ID, err)
		}
	}

	return nil
}

// scheduleTask adds a task to the cron scheduler
func (s *RecurringTaskService) scheduleTask(task models.Task) error {
	var pattern models.RecurringPattern
	if err := json.Unmarshal(task.RecurringConfig, &pattern); err != nil {
		return fmt.Errorf("invalid recurring config: %v", err)
	}

	// Get cron expression based on pattern
	cronExpr := s.getCronExpression(pattern)

	// Add to cron scheduler
	_, err := s.cron.AddFunc(cronExpr, func() {
		s.executeRecurringTask(task)
	})

	return err
}

// getCronExpression converts RecurringPattern to cron expression
func (s *RecurringTaskService) getCronExpression(pattern models.RecurringPattern) string {
	switch pattern.Type {
	case "daily":
		return fmt.Sprintf("0 0 */%d * * *", pattern.Interval)
	case "weekly":
		weekdays := ""
		for _, day := range pattern.Weekdays {
			weekdays += fmt.Sprintf(",%d", day)
		}
		return fmt.Sprintf("0 0 * * %s", weekdays[1:])
	case "monthly":
		return fmt.Sprintf("0 0 1 */%d * *", pattern.Interval)
	case "custom":
		return pattern.CronExpr
	default:
		return "0 0 * * * *" // Default to daily
	}
}

// executeRecurringTask creates a new instance of the recurring task
func (s *RecurringTaskService) executeRecurringTask(template models.Task) {
	newTask := models.Task{
		Name:         template.Name,
		Description:  template.Description,
		ScheduleTime: time.Now(),
		Priority:     template.Priority,
		Status:       "pending",
		Tags:         template.Tags,
		Metadata:     template.Metadata,
		Assignee:     template.Assignee,
		Labels:       template.Labels,
	}

	if err := s.db.Create(&newTask).Error; err != nil {
		// Log error
		return
	}

	// Create task dependencies if any
	if len(template.DependentTasks) > 0 {
		for _, depTask := range template.DependentTasks {
			depTask.ParentTaskID = &newTask.ID
			if err := s.db.Create(&depTask).Error; err != nil {
				// Log error
				continue
			}
		}
	}
}

// StopScheduler stops the cron scheduler
func (s *RecurringTaskService) StopScheduler() {
	s.cron.Stop()
}

// CreateRecurringTask creates a new recurring task
func (s *RecurringTaskService) CreateRecurringTask(task *models.Task, pattern models.RecurringPattern) error {
	// Set recurring fields
	task.IsRecurring = true
	configBytes, err := json.Marshal(pattern)
	if err != nil {
		return fmt.Errorf("failed to marshal recurring config: %v", err)
	}
	task.RecurringConfig = configBytes

	// Save to database
	if err := s.db.Create(task).Error; err != nil {
		return fmt.Errorf("failed to create recurring task: %v", err)
	}

	// Schedule the task
	return s.scheduleTask(*task)
}

// UpdateRecurringTask updates an existing recurring task
func (s *RecurringTaskService) UpdateRecurringTask(taskID uint, task *models.Task, pattern models.RecurringPattern) error {
	// Update recurring config
	configBytes, err := json.Marshal(pattern)
	if err != nil {
		return fmt.Errorf("failed to marshal recurring config: %v", err)
	}
	task.RecurringConfig = configBytes

	// Update in database
	if err := s.db.Model(&models.Task{}).Where("id = ?", taskID).Updates(task).Error; err != nil {
		return fmt.Errorf("failed to update recurring task: %v", err)
	}

	// Reschedule the task
	return s.scheduleTask(*task)
}
