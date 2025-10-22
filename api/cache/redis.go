package cache

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func ConnectRedis() *redis.Client {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379" // fallback
	}
	
	log.Printf("🔗 Подключаемся к Redis по адресу: %s", addr)
	
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})

	// Тестируем подключение
	_, err := rdb.Ping(Ctx).Result()
	if err != nil {
		log.Printf("❌ Ошибка подключения к Redis: %v", err)
		return nil
	}

	log.Println("✅ Успешно подключились к Redis")
	return rdb
}