package repository

import (
	"encoding/json"
	"time"
    "golang.org/x/crypto/bcrypt"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"github.com/tovma/ruslanzadacha/api/models"
	"github.com/tovma/ruslanzadacha/api/queue"
	"gorm.io/gorm"
)

type NotificationRepository struct {
    DB    *gorm.DB
    Cache *redis.Client
    MQCh  *amqp091.Channel 
}


func (r *NotificationRepository) CreateNotification(n *models.Notification) error {
    n.CreatedAt = time.Now()
    if err := r.DB.Create(n).Error; err != nil {
        return err
    }

    
    if r.MQCh != nil {
        data, err := json.Marshal(n)
        if err == nil {
            queue.PublishMessage(r.MQCh, "notifications", data)
        }
    }

    return nil
}
func (r *NotificationRepository) CreateUser(u *models.RegUser) (uint,error)  {
    
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
    user := models.User{
    Username:     u.Username,   
    Email:        u.Email,      
    PasswordHash: string(hashedPassword), 
   
}
    if err := r.DB.Create(&user).Error; err != nil {
        return 0,err
    }


    return user.ID,nil
}


