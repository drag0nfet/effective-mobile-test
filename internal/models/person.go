package models

import (
	"gorm.io/gorm"
)

type Person struct {
	gorm.Model
	Name        string `json:"name" gorm:"not null"`
	Surname     string `json:"surname" gorm:"not null"`
	Patronymic  string `json:"patronymic"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}
