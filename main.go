package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/task-schedulart/config"
	"github.com/task-schedulart/models"
	"github.com/task-schedulart/services"
	"go.uber.org/zap"
)

// PaginationQuery represents query parameters for pagination
type PaginationQuery struct {
	Page     int `form:"page,default=1" binding:"min=1"`
	PageSize int `form:"page_size,default=10" binding:"min=1,max=100"`
}

// TaskQuery represents query parameters for task filtering
type TaskQuery struct {
	Status   string   `form:"status"`
	Priority string   `form:"priority"`
	Tags     []string `form:"tags"`
	Search   string   `form:"search"`
	SortBy   string   `form:"sort_by,default=created_at"`
	Order    string   `form:"order,default=desc"`
	PaginationQuery
}

// convertToUint converts string ID to uint and handles errors
func convertToUint(id string) (uint, error) {
	num, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid ID format: %v", err)
	}
	return uint(num), nil
}

func main() {
	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Initialize database
	db, err := config.InitDB()
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}

	// Initialize services
	taskService := services.NewTaskService(db)
	metricsService := services.NewMetricsService()
	wsService := services.NewWebSocketService(logger)
	recurringService := services.NewRecurringTaskService(db)

	// Start WebSocket service
	go wsService.Start()

	// Start recurring task service
	go recurringService.StartScheduler()

	// Create Gin router
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// Serve static files
	r.Static("/static", "./static")
	r.LoadHTMLGlob("frontend/*.html")

	// Serve frontend
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// API Routes
	api := r.Group("/api")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "healthy"})
		})

		// Metrics endpoint
		api.GET("/metrics", gin.WrapH(metricsService.Handler()))

		// Task routes
		tasks := api.Group("/tasks")
		{
			// List tasks with filtering and pagination
			tasks.GET("", func(c *gin.Context) {
				var query TaskQuery
				if err := c.ShouldBindQuery(&query); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				tasks, total, err := taskService.GetTasksWithPagination(query.Status, query.Priority, query.Tags,
					query.Search, query.SortBy, query.Order, query.Page, query.PageSize)
				if err != nil {
					logger.Error("Failed to fetch tasks", zap.Error(err))
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"tasks": tasks,
					"pagination": gin.H{
						"current_page": query.Page,
						"page_size":    query.PageSize,
						"total_items":  total,
						"total_pages":  (total + int64(query.PageSize) - 1) / int64(query.PageSize),
					},
				})
			})

			// Create new task
			tasks.POST("", func(c *gin.Context) {
				var task models.Task
				if err := c.ShouldBindJSON(&task); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				// Set default values
				task.Status = "pending"
				task.CreatedAt = time.Now()
				task.UpdatedAt = time.Now()

				if err := taskService.CreateTask(&task); err != nil {
					logger.Error("Failed to create task", zap.Error(err))
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				// Record metrics
				metricsService.RecordTaskCreation()

				// Broadcast WebSocket update
				wsService.BroadcastTaskUpdate(services.TaskCreatedEvent, task)

				c.JSON(http.StatusCreated, task)
			})

			// Get task by ID
			tasks.GET("/:id", func(c *gin.Context) {
				taskID, err := convertToUint(c.Param("id"))
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				task, err := taskService.GetTaskByID(taskID)
				if err != nil {
					logger.Error("Failed to get task", zap.Error(err))
					c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
					return
				}

				c.JSON(http.StatusOK, task)
			})

			// Update task
			tasks.PUT("/:id", func(c *gin.Context) {
				taskID, err := convertToUint(c.Param("id"))
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				var task models.Task
				if err := c.ShouldBindJSON(&task); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				task.ID = taskID
				task.UpdatedAt = time.Now()

				if err := taskService.UpdateTask(&task); err != nil {
					logger.Error("Failed to update task", zap.Error(err))
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				// Broadcast WebSocket update
				wsService.BroadcastTaskUpdate(services.TaskUpdatedEvent, task)

				c.JSON(http.StatusOK, task)
			})

			// Update task status
			tasks.PUT("/:id/status", func(c *gin.Context) {
				taskID, err := convertToUint(c.Param("id"))
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				var req struct {
					Status string `json:"status" binding:"required,oneof=pending running completed failed"`
				}

				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				if err := taskService.UpdateTaskStatus(taskID, req.Status); err != nil {
					logger.Error("Failed to update task status", zap.Error(err))
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				// Update metrics based on status
				switch req.Status {
				case "completed":
					metricsService.RecordTaskCompletion()
				case "failed":
					metricsService.RecordTaskFailure()
				}

				// Broadcast WebSocket update
				wsService.BroadcastTaskUpdate(services.TaskStatusEvent, gin.H{
					"id":     taskID,
					"status": req.Status,
				})

				c.JSON(http.StatusOK, gin.H{"message": "Task status updated"})
			})

			// Retry failed task
			tasks.POST("/:id/retry", func(c *gin.Context) {
				taskID, err := convertToUint(c.Param("id"))
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				if err := taskService.RetryFailedTask(taskID); err != nil {
					logger.Error("Failed to retry task", zap.Error(err))
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				// Broadcast WebSocket update
				wsService.BroadcastTaskUpdate(services.TaskStatusEvent, gin.H{
					"id":     taskID,
					"status": "pending",
				})

				c.JSON(http.StatusOK, gin.H{"message": "Task scheduled for retry"})
			})

			// Delete task
			tasks.DELETE("/:id", func(c *gin.Context) {
				taskID, err := convertToUint(c.Param("id"))
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				if err := taskService.DeleteTask(taskID); err != nil {
					logger.Error("Failed to delete task", zap.Error(err))
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				// Broadcast WebSocket update
				wsService.BroadcastTaskUpdate(services.TaskDeletedEvent, gin.H{"id": taskID})

				c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
			})

			// Get tasks by tags
			tasks.GET("/tags", func(c *gin.Context) {
				tags := c.QueryArray("tags")
				if len(tags) == 0 {
					c.JSON(http.StatusBadRequest, gin.H{"error": "No tags provided"})
					return
				}

				tasks, err := taskService.GetTasksByTags(tags)
				if err != nil {
					logger.Error("Failed to fetch tasks by tags", zap.Error(err))
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusOK, tasks)
			})

			// Create recurring task
			tasks.POST("/recurring", func(c *gin.Context) {
				var task models.Task
				if err := c.ShouldBindJSON(&task); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				var pattern models.RecurringPattern
				if err := c.ShouldBindJSON(&pattern); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				if err := recurringService.CreateRecurringTask(&task, pattern); err != nil {
					logger.Error("Failed to create recurring task", zap.Error(err))
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusCreated, task)
			})
		}
	}

	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		logger.Info("No PORT environment variable, using default", zap.String("port", port))
	}

	// Start the server
	logger.Info(fmt.Sprintf("Starting server on port %s", port))
	if err := r.Run(":" + port); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
