package main

import (
	"log"


	"github.com/gin-gonic/gin"

	"github.com/tovma/ruslanzadacha/api/cache"
	"github.com/tovma/ruslanzadacha/api/db"
	"github.com/tovma/ruslanzadacha/api/queue"
	"github.com/tovma/ruslanzadacha/api/repository"
	"github.com/tovma/ruslanzadacha/api/routes"
)

func main() {
    r := gin.Default()

    dbConn := db.ConnectDb()
    redisClient := cache.ConnectRedis()


    rabbitConn, rabbitCh, err := queue.ConnectRabbitMQ()
    if err != nil {
        log.Fatalf("❌ Ошибка подключения к RabbitMQ: %v", err)
    }
    
    defer rabbitConn.Close()
    defer rabbitCh.Close()

    queue.DeclareQueue(rabbitCh, "notifications")
    
    repo := &repository.NotificationRepository{
        DB:    dbConn,
        Cache: redisClient,
        MQCh:  rabbitCh,
    }

    routes.SetupRotes(r, repo)

    log.Println("🌐 Сервер запускается на :8080")
    r.Run("0.0.0.0:8080")
}