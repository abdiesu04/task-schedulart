package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Task struct {
	Name         string `json:"name"`
	ScheduleTime string `json:"scheduleTime"`
	Priority     string `json:"priority"`
}

var tasks []Task

func main() {
	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Create Gin router
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Origin", "Content-Type"}
	r.Use(cors.New(config))

	// Basic health check endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Task Schedulart API is running",
		})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	// Task endpoints
	r.GET("/tasks", func(c *gin.Context) {
		logger.Info("Fetching tasks")
		c.JSON(http.StatusOK, tasks)
	})

	r.POST("/tasks", func(c *gin.Context) {
		logger.Info("Received new task request")
		var newTask Task
		if err := c.BindJSON(&newTask); err != nil {
			logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		tasks = append(tasks, newTask)
		logger.Info("Task created successfully")
		c.JSON(http.StatusCreated, newTask)
	})

	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		logger.Info("No PORT environment variable, using default", zap.String("port", port))
	}

	// Log the port we're using
	logger.Info(fmt.Sprintf("Starting server on port %s", port))

	// Start the server
	if err := r.Run(":" + port); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
