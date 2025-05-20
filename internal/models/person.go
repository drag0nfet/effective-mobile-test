package models

import (
	"gorm.io/gorm"
)

// Person represents a person entity
// @Description Person entity with personal details
type Person struct {
	gorm.Model
	// @Description Person's first name
	Name string `json:"name" gorm:"not null" example:"Dmitriy"`
	// @Description Person's surname
	Surname string `json:"surname" gorm:"not null" example:"Ushakov"`
	// @Description Person's patronymic (optional)
	Patronymic string `json:"patronymic" example:"Vasilevich"`
	// @Description Person's age, enriched via API
	Age int `json:"age" example:"34"`
	// @Description Person's gender, enriched via API
	Gender string `json:"gender" example:"male"`
	// @Description Person's nationality, enriched via API
	Nationality string `json:"nationality" example:"RU"`
}
