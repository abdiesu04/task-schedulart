package services

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MetricsService struct {
	TasksCreated    prometheus.Counter
	TasksCompleted  prometheus.Counter
	TasksFailed     prometheus.Counter
	TasksProcessing prometheus.Gauge
	TaskDuration    prometheus.Histogram
	TasksByStatus   *prometheus.GaugeVec
	TasksByPriority *prometheus.GaugeVec
}

func NewMetricsService() *MetricsService {
	return &MetricsService{
		TasksCreated: promauto.NewCounter(prometheus.CounterOpts{
			Name: "task_schedulart_tasks_created_total",
			Help: "The total number of created tasks",
		}),
		TasksCompleted: promauto.NewCounter(prometheus.CounterOpts{
			Name: "task_schedulart_tasks_completed_total",
			Help: "The total number of completed tasks",
		}),
		TasksFailed: promauto.NewCounter(prometheus.CounterOpts{
			Name: "task_schedulart_tasks_failed_total",
			Help: "The total number of failed tasks",
		}),
		TasksProcessing: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "task_schedulart_tasks_processing",
			Help: "The number of tasks currently being processed",
		}),
		TaskDuration: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "task_schedulart_task_duration_seconds",
			Help:    "Task execution duration in seconds",
			Buckets: prometheus.DefBuckets,
		}),
		TasksByStatus: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "task_schedulart_tasks_by_status",
				Help: "Number of tasks by status",
			},
			[]string{"status"},
		),
		TasksByPriority: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "task_schedulart_tasks_by_priority",
				Help: "Number of tasks by priority",
			},
			[]string{"priority"},
		),
	}
}

// RecordTaskCreation increments the tasks created counter
func (m *MetricsService) RecordTaskCreation() {
	m.TasksCreated.Inc()
}

// RecordTaskCompletion increments the tasks completed counter
func (m *MetricsService) RecordTaskCompletion() {
	m.TasksCompleted.Inc()
}

// RecordTaskFailure increments the tasks failed counter
func (m *MetricsService) RecordTaskFailure() {
	m.TasksFailed.Inc()
}

// UpdateTasksProcessing sets the number of tasks currently being processed
func (m *MetricsService) UpdateTasksProcessing(count float64) {
	m.TasksProcessing.Set(count)
}

// ObserveTaskDuration records the duration of a task execution
func (m *MetricsService) ObserveTaskDuration(durationSeconds float64) {
	m.TaskDuration.Observe(durationSeconds)
}

// UpdateTaskStatusMetric updates the count for a specific task status
func (m *MetricsService) UpdateTaskStatusMetric(status string, count float64) {
	m.TasksByStatus.WithLabelValues(status).Set(count)
}

// UpdateTaskPriorityMetric updates the count for a specific task priority
func (m *MetricsService) UpdateTaskPriorityMetric(priority string, count float64) {
	m.TasksByPriority.WithLabelValues(priority).Set(count)
}
