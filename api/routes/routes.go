package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tovma/ruslanzadacha/api/handlers"
	"github.com/tovma/ruslanzadacha/api/repository"
)

func SetupRotes(r *gin.Engine, repo *repository.NotificationRepository) {
	r.POST("/notifications", handlers.CreateNotificationHandler(repo))
	r.GET("/notifications", handlers.GetNotificationsHandler(repo))
	r.DELETE("/notifications/:id", handlers.DeleteNotificationsHandler(repo))
	r.GET("/notifications/:id", handlers.GetNotificationByIdHandler(repo))
}

