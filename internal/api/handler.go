package api

import (
	"github.com/drag0nfet/effective-mobile-test/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

// SetupRoutes настраивает маршруты данного API
func SetupRoutes(r *gin.Engine, personService *service.PersonService) {
	r.Use(loggingMiddleware)

	personHandler := NewPersonHandler(personService)

	// Маршруты для работы с Person
	r.POST("/persons", personHandler.CreatePerson)
	r.GET("/persons", personHandler.GetPeople)
	r.PUT("/persons/:id", personHandler.UpdatePerson)
	r.DELETE("/persons/:id", personHandler.DeletePerson)
}

// loggingMiddleware логирует входящие запросы
func loggingMiddleware(c *gin.Context) {
	start := time.Now()
	method := c.Request.Method
	path := c.Request.URL.Path

	// Выполняем запрос
	c.Next()

	// Логируем после выполнения запроса
	latency := time.Since(start)
	status := c.Writer.Status()
	logrus.WithFields(logrus.Fields{
		"method":  method,
		"path":    path,
		"status":  status,
		"latency": latency,
	}).Info("HTTP request processed")
}
