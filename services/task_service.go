package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/task-schedulart/models"
	"gorm.io/gorm"
)

type TaskService struct {
	db *gorm.DB
}

func NewTaskService(db *gorm.DB) *TaskService {
	return &TaskService{db: db}
}

// CreateTask creates a new task
func (s *TaskService) CreateTask(task *models.Task) error {
	return s.db.Create(task).Error
}

// GetTasks returns all tasks with optional filters
func (s *TaskService) GetTasks(status, priority string) ([]models.Task, error) {
	var tasks []models.Task
	query := s.db.Model(&models.Task{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if priority != "" {
		query = query.Where("priority = ?", priority)
	}

	err := query.Order("schedule_time asc").Find(&tasks).Error
	return tasks, err
}

// UpdateTaskStatus updates the status of a task
func (s *TaskService) UpdateTaskStatus(taskID uint, status string) error {
	return s.db.Model(&models.Task{}).Where("id = ?", taskID).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

// GetPendingTasks returns tasks that are scheduled to run and are pending
func (s *TaskService) GetPendingTasks() ([]models.Task, error) {
	var tasks []models.Task
	err := s.db.Where("status = ? AND schedule_time <= ?", "pending", time.Now()).
		Order("priority desc, schedule_time asc").
		Find(&tasks).Error
	return tasks, err
}

// RetryFailedTask attempts to retry a failed task
func (s *TaskService) RetryFailedTask(taskID uint) error {
	var task models.Task
	if err := s.db.First(&task, taskID).Error; err != nil {
		return err
	}

	if task.Status != "failed" {
		return errors.New("only failed tasks can be retried")
	}

	if task.RetryCount >= 3 {
		return errors.New("maximum retry attempts reached")
	}

	return s.db.Model(&task).Updates(map[string]interface{}{
		"status":      "pending",
		"retry_count": task.RetryCount + 1,
		"last_error":  "",
	}).Error
}

// DeleteTask soft deletes a task
func (s *TaskService) DeleteTask(taskID uint) error {
	return s.db.Delete(&models.Task{}, taskID).Error
}

// GetTasksByTags returns tasks with specific tags
func (s *TaskService) GetTasksByTags(tags []string) ([]models.Task, error) {
	var tasks []models.Task
	err := s.db.Where("tags && ?", tags).Find(&tasks).Error
	return tasks, err
}

// GetTasksWithPagination returns paginated tasks with filters
func (s *TaskService) GetTasksWithPagination(status, priority string, tags []string, search, sortBy, order string, page, pageSize int) ([]models.Task, int64, error) {
	var tasks []models.Task
	var total int64
	offset := (page - 1) * pageSize

	query := s.db.Model(&models.Task{})

	// Apply filters
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if priority != "" {
		query = query.Where("priority = ?", priority)
	}
	if len(tags) > 0 {
		query = query.Where("tags && ?", tags)
	}
	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Get total count before pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting 
	if sortBy != "" {
		if order != "asc" {
			order = "desc"
		}
		query = query.Order(fmt.Sprintf("%s %s", sortBy, order))
	}

	// Apply pagination
	err := query.Offset(offset).Limit(pageSize).Find(&tasks).Error
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// GetTaskByID retrieves a task by its ID
func (s *TaskService) GetTaskByID(taskID uint) (*models.Task, error) {
	var task models.Task
	if err := s.db.First(&task, taskID).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

// UpdateTask updates an existing task
func (s *TaskService) UpdateTask(task *models.Task) error {
	if task.ID == 0 {
		return errors.New("task ID is required")
	}

	// Ensure the task exists
	var existingTask models.Task
	if err := s.db.First(&existingTask, task.ID).Error; err != nil {
		return err
	}

	// Update the task
	task.UpdatedAt = time.Now()
	return s.db.Model(task).Updates(task).Error
}
