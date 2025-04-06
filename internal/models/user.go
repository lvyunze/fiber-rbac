package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex"`
	Password string
	Roles    []Role `gorm:"many2many:user_roles;"`
}
