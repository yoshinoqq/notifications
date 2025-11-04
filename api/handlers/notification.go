package handlers

import (
	"encoding/json"
	"log"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/tovma/ruslanzadacha/api/cache"
	"github.com/tovma/ruslanzadacha/api/models"
	"github.com/tovma/ruslanzadacha/api/repository"
)
func RegUserHandler(repo *repository.NotificationRepository) gin.HandlerFunc {
    return func(c *gin.Context){
           log.Println("✅ RegUserHandler вызван!")
        var n models.RegUser
        
        if err := c.BindJSON(&n); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
            return
        }
        

        userID, err := repo.CreateUser(&n)
        if err != nil {
            log.Printf("❌ Ошибка CreateUser: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save"})
            return
        }
         log.Printf("✅ Пользователь создан с ID: %d", userID)
        response := models.UserResponse{
            ID:        userID,
            Username:  n.Username,
            Email:     n.Email,
            CreatedAt: time.Now(),
        }
        c.JSON(http.StatusCreated, response) 
    }
}
func LoginUserHandler(repo *repository.NotificationRepository) gin.HandlerFunc {
    return func(c *gin.Context){
        
    }
}
func CreateNotificationHandler(repo *repository.NotificationRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        var n models.Notification
        if err := c.BindJSON(&n); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
            return
        }

        if err := repo.CreateNotification(&n); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save"})
            return
        }

        c.JSON(http.StatusOK, n)
    }
}

func GetNotificationsHandler(repo *repository.NotificationRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        const cacheKey = "notifications"

        // Проверяем кеш Redis
        cached, err := repo.Cache.Get(cache.Ctx, cacheKey).Result()
        if err == nil && cached != "" {
            var notifications []models.Notification
            if err := json.Unmarshal([]byte(cached), &notifications); err != nil {
                log.Println("Ошибка при парсинге кеша:", err)
            } else {
                log.Println(" Получено из Redis кеша:", cacheKey)
                c.JSON(http.StatusOK, gin.H{
                    "source": "cache",
                    "data":   notifications,
                })
                return
            }
        }


        var notifications []models.Notification
        if err := repo.DB.Find(&notifications).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch notifications"})
            return
        }

        // Сохраняем в кеш на 5 минут
        data, err := json.Marshal(notifications)
        if err != nil {
            log.Println("❌ Ошибка при сериализации для кеша:", err)
        } else {
            repo.Cache.Set(cache.Ctx, cacheKey, data, 5*time.Minute)
        }

        log.Println("✅ Получено из базы данных")
        c.JSON(http.StatusOK, gin.H{
            "source": "database",
            "data":   notifications,
        })
    }
}
func DeleteNotificationsHandler(repo *repository.NotificationRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        
		id := c.Param("id")
        
        if err := repo.DB.Delete(&models.Notification{}, id).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete notifications"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "notification deleted"} )
    }
}
func GetNotificationByIdHandler(repo *repository.NotificationRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
		var n models.Notification
        id := c.Param("id")
        
        if err := repo.DB.First(&n,id).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch notification"})
            return
        }

        c.JSON(http.StatusOK,n )
    }
}