package api

import (
	"github.com/drag0nfet/effective-mobile-test/internal/models"
	"github.com/drag0nfet/effective-mobile-test/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type PersonHandler struct {
	service *service.PersonService
}

func NewPersonHandler(service *service.PersonService) *PersonHandler {
	return &PersonHandler{service: service}
}

func SetupRoutes(r *gin.Engine, service *service.PersonService) {
	handler := NewPersonHandler(service)
	r.POST("/persons", handler.CreatePerson)
	r.GET("/persons", handler.GetPeople)
	r.PUT("/persons/:id", handler.UpdatePerson)
	r.DELETE("/persons/:id", handler.DeletePerson)
}

func (h *PersonHandler) CreatePerson(c *gin.Context) {
	var person models.Person
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreatePerson(&person); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, person)
}

// GetPeople godoc
// @Summary Get list of people
// @Description Retrieve a list of people with optional filters and pagination
// @Accept json
// @Produce json
// @Param name query string false "Filter by name"
// @Param surname query string false "Filter by surname"
// @Param patronymic query string false "Filter by patronymic"
// @Param age query int false "Filter by age"
// @Param gender query string false "Filter by gender"
// @Param nationality query string false "Filter by nationality"
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Items per page (default 10)"
// @Success 200 {array} models.Person
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /persons [get]
func (h *PersonHandler) GetPeople(c *gin.Context) {
	// Получение параметров фильтрации
	filters := make(map[string]interface{})
	if name := c.Query("name"); name != "" {
		filters["name"] = name
	}
	if surname := c.Query("surname"); surname != "" {
		filters["surname"] = surname
	}
	if patronymic := c.Query("patronymic"); patronymic != "" {
		filters["patronymic"] = patronymic
	}
	if age := c.Query("age"); age != "" {
		if ageVal, err := strconv.Atoi(age); err == nil {
			filters["age"] = ageVal
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid age format"})
			return
		}
	}
	if gender := c.Query("gender"); gender != "" {
		filters["gender"] = gender
	}
	if nationality := c.Query("nationality"); nationality != "" {
		filters["nationality"] = nationality
	}

	// Получение параметров пагинации
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Вызов сервиса для получения списка людей
	people, err := h.service.GetPeople(filters, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, people)
}

// UpdatePerson godoc
// @Summary Update a person
// @Description Update an existing person by ID
// @Accept json
// @Produce json
// @Param id path int true "Person ID"
// @Param person body models.Person true "Updated person data"
// @Success 200 {object} models.Person
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /persons/{id} [put]
func (h *PersonHandler) UpdatePerson(c *gin.Context) {
	// Получение ID из пути
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	var person models.Person
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	person.ID = uint(id)

	// Вызов сервиса для обновления
	if err := h.service.UpdatePerson(&person); err != nil {
		if err.Error() == "person not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "person not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, person)
}

// DeletePerson godoc
// @Summary Delete a person
// @Description Delete a person by ID
// @Produce json
// @Param id path int true "Person ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /persons/{id} [delete]
func (h *PersonHandler) DeletePerson(c *gin.Context) {
	// Получение ID из пути
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	// Вызов сервиса для удаления
	if err := h.service.DeletePerson(uint(id)); err != nil {
		if err.Error() == "person not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "person not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
