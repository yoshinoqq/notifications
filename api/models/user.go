package models
import "gorm.io/gorm"
import "time"

type User struct {
    gorm.Model                   
    Username     string `json:"username" gorm:"unique;not null"`
    Email        string `json:"email" gorm:"unique;not null"`
    PasswordHash string `json:"-" gorm:"not null"`
    
}
type RegUser struct {
    Username string `json:"username" binding:"required, min=3"`
    Email    string `json:"email" binding:"required,email"` 
    Password string `json:"password" binding:"required,min=6"`
}
type UserResponse struct {
    ID        uint      `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}