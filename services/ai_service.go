package services

import (
	"math"
	"sort"
	"time"

	"github.com/task-schedulart/models"
	"gorm.io/gorm"
)

type AIService struct {
	db *gorm.DB
}

func NewAIService(db *gorm.DB) *AIService {
	return &AIService{db: db}
}

// TaskAnalytics represents analytics data for tasks
type TaskAnalytics struct {
	CompletionRate       float64            `json:"completionRate"`
	AverageCompletion    float64            `json:"averageCompletion"`
	ProductivityHours    map[int]float64    `json:"productivityHours"`
	TagPerformance       map[string]float64 `json:"tagPerformance"`
	PriorityDistribution map[string]int     `json:"priorityDistribution"`
	TrendAnalysis        []TrendPoint       `json:"trendAnalysis"`
}

type TrendPoint struct {
	Date  time.Time `json:"date"`
	Value float64   `json:"value"`
}

// OptimizeTaskSchedule uses AI to optimize task scheduling
func (s *AIService) OptimizeTaskSchedule() error {
	var tasks []models.Task
	if err := s.db.Where("status = ?", "pending").Find(&tasks).Error; err != nil {
		return err
	}

	// Calculate task scores and sort by priority
	type TaskScore struct {
		Task  models.Task
		Score float64
	}

	var scoredTasks []TaskScore
	for _, task := range tasks {
		score := s.calculateTaskScore(task)
		scoredTasks = append(scoredTasks, TaskScore{Task: task, Score: score})
	}

	// Sort tasks by score
	sort.Slice(scoredTasks, func(i, j int) bool {
		return scoredTasks[i].Score > scoredTasks[j].Score
	})

	// Update task schedule times based on optimization
	for i, st := range scoredTasks {
		optimalTime := s.calculateOptimalTime(st.Task, i)
		st.Task.ScheduleTime = optimalTime
		st.Task.Priority_Score = st.Score
		if err := s.db.Save(&st.Task).Error; err != nil {
			return err
		}
	}

	return nil
}

// calculateTaskScore uses multiple factors to determine task priority
func (s *AIService) calculateTaskScore(task models.Task) float64 {
	var score float64

	// Priority weight
	priorityWeight := map[string]float64{
		"high":   1.0,
		"medium": 0.6,
		"low":    0.3,
	}
	score += priorityWeight[task.Priority]

	// Due date urgency
	if task.DueDate != nil {
		timeUntilDue := task.DueDate.Sub(time.Now())
		urgencyScore := 1.0 - (timeUntilDue.Hours() / (24.0 * 7)) // Higher score for closer deadlines
		score += math.Max(0, urgencyScore)
	}

	// Historical performance with similar tasks
	historicalScore := s.analyzeHistoricalPerformance(task)
	score += historicalScore

	// Dependencies weight
	if len(task.DependentTasks) > 0 {
		score += 0.5 // Boost score for tasks with dependencies
	}

	return score
}

// calculateOptimalTime determines the best time to schedule a task
func (s *AIService) calculateOptimalTime(task models.Task, position int) time.Time {
	baseTime := time.Now()

	// Consider user's productive hours
	productiveHours := s.analyzeProductiveHours()

	// Add delay based on position and available resources
	delay := time.Duration(position*30) * time.Minute

	// Find next optimal time slot
	optimalTime := baseTime.Add(delay)
	for !s.isOptimalTimeSlot(optimalTime, productiveHours) {
		optimalTime = optimalTime.Add(30 * time.Minute)
	}

	return optimalTime
}

// analyzeHistoricalPerformance analyzes past task performance
func (s *AIService) analyzeHistoricalPerformance(task models.Task) float64 {
	var similarTasks []models.Task
	s.db.Where("tags && ? AND status = ?", task.Tags, "completed").Find(&similarTasks)

	if len(similarTasks) == 0 {
		return 0.5 // Default score for new task types
	}

	var successRate float64
	for _, t := range similarTasks {
		if t.ActualTime <= t.EstimatedTime {
			successRate += 1.0
		}
	}

	return successRate / float64(len(similarTasks))
}

// analyzeProductiveHours determines the most productive hours based on task completion history
func (s *AIService) analyzeProductiveHours() map[int]float64 {
	productiveHours := make(map[int]float64)
	var completedTasks []models.Task

	s.db.Where("status = ?", "completed").Find(&completedTasks)

	for _, task := range completedTasks {
		hour := task.UpdatedAt.Hour()
		if task.ActualTime > 0 && task.EstimatedTime > 0 {
			efficiency := float64(task.EstimatedTime) / float64(task.ActualTime)
			productiveHours[hour] += efficiency
		}
	}

	// Normalize scores
	var maxScore float64
	for _, score := range productiveHours {
		if score > maxScore {
			maxScore = score
		}
	}

	for hour := range productiveHours {
		productiveHours[hour] /= maxScore
	}

	return productiveHours
}

// isOptimalTimeSlot checks if a given time is optimal for task scheduling
func (s *AIService) isOptimalTimeSlot(t time.Time, productiveHours map[int]float64) bool {
	hour := t.Hour()
	score := productiveHours[hour]

	// Consider it optimal if:
	// 1. It's during productive hours (score > 0.7)
	// 2. Not too many tasks are already scheduled
	// 3. Not during typical off-hours (11 PM - 5 AM)
	if score > 0.7 && hour >= 5 && hour <= 23 {
		var conflictingTasks int64
		s.db.Model(&models.Task{}).
			Where("schedule_time BETWEEN ? AND ?",
				t.Add(-30*time.Minute),
				t.Add(30*time.Minute)).
			Count(&conflictingTasks)

		return conflictingTasks < 3
	}

	return false
}

// GenerateTaskAnalytics generates comprehensive analytics for tasks
func (s *AIService) GenerateTaskAnalytics(startDate, endDate time.Time) (*TaskAnalytics, error) {
	var tasks []models.Task
	if err := s.db.Where("created_at BETWEEN ? AND ?", startDate, endDate).Find(&tasks).Error; err != nil {
		return nil, err
	}

	analytics := &TaskAnalytics{
		ProductivityHours:    make(map[int]float64),
		TagPerformance:       make(map[string]float64),
		PriorityDistribution: make(map[string]int),
	}

	var completed, total float64
	tagCounts := make(map[string]int)

	for _, task := range tasks {
		total++
		if task.Status == "completed" {
			completed++
		}

		// Track priority distribution
		analytics.PriorityDistribution[task.Priority]++

		// Analyze tags
		for _, tag := range task.Tags {
			tagCounts[tag]++
			if task.Status == "completed" {
				analytics.TagPerformance[tag]++
			}
		}

		// Track productivity by hour
		if task.Status == "completed" {
			hour := task.UpdatedAt.Hour()
			analytics.ProductivityHours[hour]++
		}
	}

	// Calculate completion rate
	analytics.CompletionRate = (completed / total) * 100

	// Calculate average completion time
	var totalCompletionTime float64
	var completedCount float64
	for _, task := range tasks {
		if task.Status == "completed" && task.ActualTime > 0 {
			totalCompletionTime += float64(task.ActualTime)
			completedCount++
		}
	}
	if completedCount > 0 {
		analytics.AverageCompletion = totalCompletionTime / completedCount
	}
	// Normalize tag performance
	for tag, completions := range analytics.TagPerformance {
		total := float64(tagCounts[tag])
		analytics.TagPerformance[tag] = (completions / total) * 100
	}

	// Generate trend analysis
	analytics.TrendAnalysis = s.generateTrendAnalysis(startDate, endDate)

	return analytics, nil
}

// generateTrendAnalysis generates trend data for task completion over time
func (s *AIService) generateTrendAnalysis(startDate, endDate time.Time) []TrendPoint {
	var trends []TrendPoint
	interval := endDate.Sub(startDate) / 10 // Split into 10 points

	for point := startDate; point.Before(endDate); point = point.Add(interval) {
		var total, completed int64 // Changed to int64 to match Count() requirements
		s.db.Model(&models.Task{}).
			Where("created_at BETWEEN ? AND ?", point, point.Add(interval)).
			Count(&total)

		s.db.Model(&models.Task{}).
			Where("created_at BETWEEN ? AND ? AND status = ?", point, point.Add(interval), "completed").
			Count(&completed)

		var completionRate float64
		if total > 0 {
			completionRate = float64(completed) / float64(total) * 100
		}

		trends = append(trends, TrendPoint{
			Date:  point,
			Value: completionRate,
		})
	}

	return trends
}
