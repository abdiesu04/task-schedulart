package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/task-schedulart/models"
	"gorm.io/gorm"
)

type NotificationService struct {
	db *gorm.DB
}

type NotificationTemplate struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Type     string `json:"type"`
	Subject  string `json:"subject"`
	Template string `json:"template"`
}

type NotificationChannel struct {
	Type    string          `json:"type"` // email, slack, webhook, push
	Config  json.RawMessage `json:"config"`
	Enabled bool            `json:"enabled"`
}

type EmailConfig struct {
	SMTP     string `json:"smtp"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type SlackConfig struct {
	WebhookURL string `json:"webhookUrl"`
	Channel    string `json:"channel"`
}

type WebhookConfig struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
}

func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{db: db}
}

// SendTaskNotification sends notifications for task events
func (s *NotificationService) SendTaskNotification(task *models.Task, event string) error {
	// Get notification template
	var template NotificationTemplate
	if err := s.db.Where("type = ?", event).First(&template).Error; err != nil {
		return fmt.Errorf("template not found: %v", err)
	}

	// Get notification channels
	var channels []NotificationChannel
	if err := s.db.Where("enabled = ?", true).Find(&channels).Error; err != nil {
		return fmt.Errorf("failed to get channels: %v", err)
	}

	// Prepare notification data
	data := map[string]interface{}{
		"task":      task,
		"timestamp": time.Now(),
		"event":     event,
	}

	// Send to each channel
	for _, channel := range channels {
		if err := s.sendToChannel(channel, template, data); err != nil {
			// Log error but continue with other channels
			fmt.Printf("Failed to send to channel %s: %v\n", channel.Type, err)
		}
	}

	return nil
}

// sendToChannel sends a notification through a specific channel
func (s *NotificationService) sendToChannel(channel NotificationChannel, tmpl NotificationTemplate, data map[string]interface{}) error {
	// Parse and execute template
	t, err := template.New("notification").Parse(tmpl.Template)
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}

	var content bytes.Buffer
	if err := t.Execute(&content, data); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	switch channel.Type {
	case "email":
		return s.sendEmail(channel.Config, tmpl.Subject, content.String())
	case "slack":
		return s.sendSlack(channel.Config, content.String())
	case "webhook":
		return s.sendWebhook(channel.Config, content.String())
	default:
		return fmt.Errorf("unsupported channel type: %s", channel.Type)
	}
}

// sendEmail sends an email notification
func (s *NotificationService) sendEmail(configData json.RawMessage, subject, content string) error {
	var config EmailConfig
	if err := json.Unmarshal(configData, &config); err != nil {
		return fmt.Errorf("invalid email config: %v", err)
	}

	// Implement email sending logic here
	// For example, using net/smtp or a third-party email service
	return nil
}

// sendSlack sends a Slack notification
func (s *NotificationService) sendSlack(configData json.RawMessage, content string) error {
	var config SlackConfig
	if err := json.Unmarshal(configData, &config); err != nil {
		return fmt.Errorf("invalid slack config: %v", err)
	}

	payload := map[string]interface{}{
		"text":    content,
		"channel": config.Channel,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(config.WebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack API returned status: %d", resp.StatusCode)
	}

	return nil
}

// sendWebhook sends a webhook notification
func (s *NotificationService) sendWebhook(configData json.RawMessage, content string) error {
	var config WebhookConfig
	if err := json.Unmarshal(configData, &config); err != nil {
		return fmt.Errorf("invalid webhook config: %v", err)
	}

	// Create request
	req, err := http.NewRequest(config.Method, config.URL, bytes.NewBufferString(content))
	if err != nil {
		return err
	}

	// Add headers
	for key, value := range config.Headers {
		req.Header.Add(key, value)
	}

	// Send request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook returned status: %d", resp.StatusCode)
	}

	return nil
}

// CreateNotificationTemplate creates a new notification template
func (s *NotificationService) CreateNotificationTemplate(template *NotificationTemplate) error {
	return s.db.Create(template).Error
}

// UpdateNotificationTemplate updates an existing notification template
func (s *NotificationService) UpdateNotificationTemplate(template *NotificationTemplate) error {
	return s.db.Save(template).Error
}

// ConfigureChannel configures a notification channel
func (s *NotificationService) ConfigureChannel(channel *NotificationChannel) error {
	// Validate channel configuration
	switch channel.Type {
	case "email":
		var config EmailConfig
		if err := json.Unmarshal(channel.Config, &config); err != nil {
			return fmt.Errorf("invalid email configuration: %v", err)
		}
	case "slack":
		var config SlackConfig
		if err := json.Unmarshal(channel.Config, &config); err != nil {
			return fmt.Errorf("invalid slack configuration: %v", err)
		}
	case "webhook":
		var config WebhookConfig
		if err := json.Unmarshal(channel.Config, &config); err != nil {
			return fmt.Errorf("invalid webhook configuration: %v", err)
		}
	default:
		return fmt.Errorf("unsupported channel type: %s", channel.Type)
	}

	return s.db.Save(channel).Error
}
