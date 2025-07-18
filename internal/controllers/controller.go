package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"smthtozip/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	taskService    *services.TaskService
	archiveService *services.Archiver
}

func NewTaskController(taskService *services.TaskService, archiveService *services.Archiver) *TaskController {
	return &TaskController{
		taskService:    taskService,
		archiveService: archiveService,
	}
}

func (tc *TaskController) CreateTask(c *gin.Context) {
	task, err := tc.taskService.CreateTask()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"task_id": task.ID})
}

func (ctrl *TaskController) AddURL(c *gin.Context) {
	taskID := c.Param("id")

	var req struct {
		URLs []string `json:"urls"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println(req)
	if len(req.URLs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "at least one URL required"})
		return
	}

	if len(req.URLs) > 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "maximum 3 URLs per request"})
		return
	}

	taskIDint, err := strconv.Atoi(taskID)
	if err != nil {
		fmt.Println("Error converting string to int:", err)
		return
	}

	for _, url := range req.URLs {
		if err := ctrl.taskService.AddURL(taskIDint, url); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
				"url":   url,
			})
			return
		}
	}

	task, err := ctrl.taskService.GetByID(taskIDint)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	log.Println("Task URLs before processing:", task.URLs)

	if len(task.URLs) > 0 {
		go func() {
			log.Println("Starting archive process")
			if err := ctrl.archiveService.Process(task, "http://localhost:8080"); err != nil {
				log.Printf("Failed to create archive for task %d: %v", taskIDint, err)
			}
		}()
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "urls added",
		"total":       len(task.URLs),
		"task_id":     taskID,
		"archivepath": "/archives/" + strconv.Itoa(taskIDint) + ".zip", // Добавляем ссылку
	})
}

func (tc *TaskController) GetTaskStatus(c *gin.Context) {
	taskID := c.Param("id")
	taskIDint, err := strconv.Atoi(taskID)
	if err != nil {
		fmt.Println("Error converting string to int:", err)
		return
	}
	task, exists := tc.taskService.GetByID(taskIDint)
	if exists != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	response := gin.H{
		"status":     task.Status,
		"urls_count": len(task.URLs),
	}

	if task.Status == "completed" {
		response["download_url"] = task.ArchivePath // Возвращаем полный URL
	}

	c.JSON(http.StatusOK, response)

}

func (tc *TaskController) DownloadArchive(c *gin.Context) {
	filename := c.Param("filename")
	path := filepath.Join("archives", filename)

	log.Println(path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	c.FileAttachment(path, filename)
}
