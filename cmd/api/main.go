package main

import (
	"github.com/drag0nfet/effective-mobile-test/internal/api"
	"github.com/drag0nfet/effective-mobile-test/internal/config"
	"github.com/drag0nfet/effective-mobile-test/internal/repository"
	"github.com/drag0nfet/effective-mobile-test/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"os"
	"time"
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

	// Перенаправляем логи в файл, который инициализируется при запуске сервера
	timestamp := time.Now().Format("2006-01-02-15-04-05")
	logFile, err := os.OpenFile("app_"+timestamp+".log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		logrus.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	logrus.SetOutput(logFile)

	// Подключение к БД
	repo, err := repository.NewPersonRepository(cfg)
	if err != nil {
		logrus.Fatalf("Failed to initialize repository: %v", err)
	}

	// Инициализация сервиса
	svc := service.NewPersonService(repo)

	r := gin.New()
	r.Use(gin.Recovery()) // Оставляем только GIN-обработку паник, чтобы не было захламления логов

	// Настройка маршрутов
	api.SetupRoutes(r, svc)

	// Инициализация swag-запросов
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	logrus.Infof("Starting server on port %s", cfg.APIPort)
	err = r.Run(":" + cfg.APIPort)
	if err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}
