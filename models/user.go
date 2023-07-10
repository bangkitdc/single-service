package models

import uuid "github.com/satori/go.uuid"

type User struct {
	ID       uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Username string    `gorm:"varchar(255)" json:"username"`
	Name     string    `gorm:"varchar(255)" json:"name"`
	Password string    `gorm:"varchar(255)" json:"password"`
}
