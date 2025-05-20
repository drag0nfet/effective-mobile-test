package main

import (
	_ "github.com/drag0nfet/effective-mobile-test/docs"
	docs "github.com/drag0nfet/effective-mobile-test/docs"
	"github.com/drag0nfet/effective-mobile-test/internal/api"
	"github.com/drag0nfet/effective-mobile-test/internal/config"
	"github.com/drag0nfet/effective-mobile-test/internal/repository"
	"github.com/drag0nfet/effective-mobile-test/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"io"
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

	// Создаем директорию logs, где будут храниться копии логов
	if err = os.MkdirAll("logs", 0755); err != nil {
		logrus.Fatalf("Failed to create logs directory: %v", err)
	}

	// Открываем файл для логов
	timestamp := time.Now().Format("2006-01-02-15-04-05")
	logFile, err := os.OpenFile("logs/app_"+timestamp+".log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		logrus.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	// Настраиваем дублирование логов и в консоль для отладки в моменте, и в файл для архива
	mw := io.MultiWriter(os.Stdout, logFile)
	logrus.SetOutput(mw)

	// Подключение к БД
	repo, err := repository.NewPersonRepository(cfg)
	if err != nil {
		logrus.Fatalf("Failed to initialize repository: %v", err)
	}

	// Инициализация сервиса
	svc := service.NewPersonService(repo, cfg)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery()) // Оставляем только GIN-обработку паник

	// Можно только через локалхост подключиться, можно убрать или заменить на нужные маски
	err = r.SetTrustedProxies([]string{"127.0.0.1" /*, "", "", ...*/})
	if err != nil {
		logrus.Fatalf("Failed to set trusted proxies: %v", err)
	}

	// Настройка маршрутов
	api.SetupRoutes(r, svc)

	// Инициализация swag-запросов
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	docs.SwaggerInfo.Host = "localhost:" + cfg.APIPort

	logrus.Infof("Starting server on port %s", cfg.APIPort)
	err = r.Run(":" + cfg.APIPort)
	if err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}
