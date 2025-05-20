package repository

import (
	"fmt"
	"github.com/drag0nfet/effective-mobile-test/internal/config"
	"github.com/drag0nfet/effective-mobile-test/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PersonRepository struct {
	Db *gorm.DB
}

func NewPersonRepository(cfg *config.Config) (*PersonRepository, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	logrus.Debug("Connecting to DB...")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	logrus.Debug("Successfully connected to DB. Starting DB-migrations...")
	err = RunMigrations(db)
	if err != nil {
		return nil, err
	}
	logrus.Debug("DB-migrations completed")

	return &PersonRepository{Db: db}, nil
}

func (r *PersonRepository) CreatePerson(person *models.Person) error {
	return r.Db.Create(person).Error
}

func (r *PersonRepository) FindPersonByID(id uint) (*models.Person, error) {
	var person models.Person
	err := r.Db.First(&person, id).Error
	return &person, err
}

func (r *PersonRepository) UpdatePerson(person *models.Person) error {
	return r.Db.Save(person).Error
}

func (r *PersonRepository) DeletePerson(id uint) error {
	return r.Db.Delete(&models.Person{}, id).Error
}

func (r *PersonRepository) FindPeople(filters map[string]interface{}, offset, limit int) ([]models.Person, error) {
	var people []models.Person
	query := r.Db.Offset(offset).Limit(limit)
	for key, value := range filters {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}
	err := query.Find(&people).Error
	return people, err
}
