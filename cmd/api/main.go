package main

import (
	"github.com/drag0nfet/effective-mobile-test/internal/config"
	"github.com/drag0nfet/effective-mobile-test/internal/repository"
	_ "github.com/drag0nfet/effective-mobile-test/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		logrus.Fatalf("Failed to load config: %v", err)
	}

	// Инициализация логгера
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)

	// Подключение к БД
	/*repo*/
	_, err = repository.NewPersonRepository(cfg)
	if err != nil {
		logrus.Fatalf("Failed to initialize repository: %v", err)
	}

	// Инициализация Gin
	r := gin.Default()

	// Запрос-тестер, позже здесь будут основные запросы
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Инициализация swag-запросов
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	logrus.Infof("Starting server on port %s", cfg.APIPort)
	err = r.Run(":" + cfg.APIPort)
	if err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}
