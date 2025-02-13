package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/task-schedulart/models"
	"gorm.io/gorm"
)

type CollaborationService struct {
	db *gorm.DB
}

type Team struct {
	ID          uint         `json:"id" gorm:"primaryKey"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
	Members     []TeamMember `json:"members" gorm:"foreignKey:TeamID"`
}

type TeamMember struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	TeamID    uint      `json:"teamId"`
	UserID    uint      `json:"userId"`
	Role      string    `json:"role"` // admin, member, viewer
	JoinedAt  time.Time `json:"joinedAt"`
	InvitedBy uint      `json:"invitedBy"`
}

type Comment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	TaskID    uint      `json:"taskId"`
	UserID    uint      `json:"userId"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Mentions  []uint    `json:"mentions" gorm:"type:integer[]"`
}

type ActivityLog struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	TaskID    uint      `json:"taskId"`
	UserID    uint      `json:"userId"`
	Action    string    `json:"action"`
	Details   string    `json:"details"`
	Timestamp time.Time `json:"timestamp"`
}

func NewCollaborationService(db *gorm.DB) *CollaborationService {
	return &CollaborationService{db: db}
}

// CreateTeam creates a new team
func (s *CollaborationService) CreateTeam(team *Team, creatorID uint) error {
	// Start transaction
	tx := s.db.Begin()
	if err := tx.Error; err != nil {
		return err
	}

	// Create team
	if err := tx.Create(team).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Add creator as admin
	member := TeamMember{
		TeamID:    team.ID,
		UserID:    creatorID,
		Role:      "admin",
		JoinedAt:  time.Now(),
		InvitedBy: creatorID,
	}

	if err := tx.Create(&member).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// InviteToTeam invites a user to join a team
func (s *CollaborationService) InviteToTeam(teamID, inviterID, inviteeID uint, role string) error {
	// Verify inviter has permission
	var inviter TeamMember
	if err := s.db.Where("team_id = ? AND user_id = ? AND role = ?", teamID, inviterID, "admin").
		First(&inviter).Error; err != nil {
		return fmt.Errorf("unauthorized: only admins can invite members")
	}

	// Create new member
	member := TeamMember{
		TeamID:    teamID,
		UserID:    inviteeID,
		Role:      role,
		JoinedAt:  time.Now(),
		InvitedBy: inviterID,
	}

	return s.db.Create(&member).Error
}

// AddComment adds a comment to a task
func (s *CollaborationService) AddComment(comment *Comment) error {
	// Extract mentions from content and store them
	// This is a simple implementation - in practice, you'd want more sophisticated mention detection
	if err := s.db.Create(comment).Error; err != nil {
		return err
	}

	// Log activity
	activity := ActivityLog{
		TaskID:    comment.TaskID,
		UserID:    comment.UserID,
		Action:    "comment_added",
		Details:   fmt.Sprintf("Comment added: %s", comment.Content),
		Timestamp: time.Now(),
	}

	return s.db.Create(&activity).Error
}

// GetTaskComments retrieves all comments for a task
func (s *CollaborationService) GetTaskComments(taskID uint) ([]Comment, error) {
	var comments []Comment
	err := s.db.Where("task_id = ?", taskID).
		Order("created_at desc").
		Find(&comments).Error
	return comments, err
}

// GetTaskActivity retrieves activity history for a task
func (s *CollaborationService) GetTaskActivity(taskID uint) ([]ActivityLog, error) {
	var activities []ActivityLog
	err := s.db.Where("task_id = ?", taskID).
		Order("timestamp desc").
		Find(&activities).Error
	return activities, err
}

// ShareTask shares a task with another team
func (s *CollaborationService) ShareTask(taskID, fromTeamID, toTeamID uint, sharerID uint) error {
	// Verify sharer has permission
	var member TeamMember
	if err := s.db.Where("team_id = ? AND user_id = ?", fromTeamID, sharerID).
		First(&member).Error; err != nil {
		return fmt.Errorf("unauthorized: user not a member of source team")
	}

	// Get original task
	var task models.Task
	if err := s.db.First(&task, taskID).Error; err != nil {
		return err
	}

	// Create shared task
	sharedTask := task
	sharedTask.ID = 0 // Clear ID for new record
	sharedTask.ParentTaskID = &taskID

	// Add sharing metadata
	metadata := map[string]interface{}{
		"shared_from": fromTeamID,
		"shared_by":   sharerID,
		"shared_at":   time.Now(),
	}
	metadataBytes, _ := json.Marshal(metadata)
	sharedTask.Metadata = string(metadataBytes)

	return s.db.Create(&sharedTask).Error
}

// GetTeamTasks retrieves all tasks for a team
func (s *CollaborationService) GetTeamTasks(teamID uint) ([]models.Task, error) {
	var tasks []models.Task
	err := s.db.Where("team_id = ?", teamID).
		Order("created_at desc").
		Find(&tasks).Error
	return tasks, err
}

// GetTeamMembers retrieves all members of a team
func (s *CollaborationService) GetTeamMembers(teamID uint) ([]TeamMember, error) {
	var members []TeamMember
	err := s.db.Where("team_id = ?", teamID).
		Find(&members).Error
	return members, err
}

// UpdateMemberRole updates a team member's role
func (s *CollaborationService) UpdateMemberRole(teamID, userID uint, newRole string, updaterID uint) error {
	// Verify updater has admin permission
	var updater TeamMember
	if err := s.db.Where("team_id = ? AND user_id = ? AND role = ?", teamID, updaterID, "admin").
		First(&updater).Error; err != nil {
		return fmt.Errorf("unauthorized: only admins can update roles")
	}

	return s.db.Model(&TeamMember{}).
		Where("team_id = ? AND user_id = ?", teamID, userID).
		Update("role", newRole).Error
}

// RemoveFromTeam removes a member from a team
func (s *CollaborationService) RemoveFromTeam(teamID, userID, removerID uint) error {
	// Verify remover has admin permission
	var remover TeamMember
	if err := s.db.Where("team_id = ? AND user_id = ? AND role = ?", teamID, removerID, "admin").
		First(&remover).Error; err != nil {
		return fmt.Errorf("unauthorized: only admins can remove members")
	}

	return s.db.Where("team_id = ? AND user_id = ?", teamID, userID).
		Delete(&TeamMember{}).Error
}

// LogActivity logs an activity
func (s *CollaborationService) LogActivity(activity *ActivityLog) error {
	return s.db.Create(activity).Error
}
