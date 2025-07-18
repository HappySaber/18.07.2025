package routes

import (
	"smthtozip/internal/controllers"
	"smthtozip/internal/services"

	"github.com/gin-gonic/gin"
)

func Routes(c *gin.Engine) {
	fileStorage := services.NewFileStorage("./archives")
	archiveService := services.NewArchiver(*fileStorage)
	taskService := services.NewTaskService(archiveService)
	ctrl := controllers.NewTaskController(taskService, archiveService)
	api := c.Group("")
	{
		api.POST("/tasks", ctrl.CreateTask)
		api.POST("/tasks/:id/urls", ctrl.AddURL)
		api.GET("/tasks/:id", ctrl.GetTaskStatus)
		//api.GET("/archives/:filename", ctrl.DownloadArchive)
	}
}
