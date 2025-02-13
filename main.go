package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/task-schedulart/config"
	"github.com/task-schedulart/models"
	"github.com/task-schedulart/services"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Initialize database
	db, err := config.InitDB()
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}

	// Initialize task service
	taskService := services.NewTaskService(db)

	// Create Gin router
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Origin", "Content-Type"}
	r.Use(cors.New(config))

	// API Routes
	api := r.Group("/api")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "healthy"})
		})

		// Task routes
		tasks := api.Group("/tasks")
		{
			// Get all tasks with optional filters
			tasks.GET("", func(c *gin.Context) {
				status := c.Query("status")
				priority := c.Query("priority")

				tasks, err := taskService.GetTasks(status, priority)
				if err != nil {
					logger.Error("Failed to fetch tasks", zap.Error(err))
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, tasks)
			})

			// Create new task
			tasks.POST("", func(c *gin.Context) {
				var task models.Task
				if err := c.ShouldBindJSON(&task); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				if err := taskService.CreateTask(&task); err != nil {
					logger.Error("Failed to create task", zap.Error(err))
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusCreated, task)
			})

			// Update task status
			tasks.PUT("/:id/status", func(c *gin.Context) {
				taskID := c.Param("id")
				var req struct {
					Status string `json:"status"`
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

				c.JSON(http.StatusOK, gin.H{"message": "Task status updated"})
			})

			// Retry failed task
			tasks.POST("/:id/retry", func(c *gin.Context) {
				taskID := c.Param("id")
				if err := taskService.RetryFailedTask(taskID); err != nil {
					logger.Error("Failed to retry task", zap.Error(err))
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusOK, gin.H{"message": "Task scheduled for retry"})
			})

			// Delete task
			tasks.DELETE("/:id", func(c *gin.Context) {
				taskID := c.Param("id")
				if err := taskService.DeleteTask(taskID); err != nil {
					logger.Error("Failed to delete task", zap.Error(err))
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
			})

			// Get tasks by tags
			tasks.GET("/tags", func(c *gin.Context) {
				tags := c.QueryArray("tags")
				tasks, err := taskService.GetTasksByTags(tags)
				if err != nil {
					logger.Error("Failed to fetch tasks by tags", zap.Error(err))
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusOK, tasks)
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
