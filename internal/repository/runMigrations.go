package repository

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	database, err := db.DB()
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(database, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/migrations",
		"postgres", driver)
	if err != nil {
		return err
	}
	return m.Up()
}
