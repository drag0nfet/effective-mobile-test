package service

import (
	"github.com/drag0nfet/effective-mobile-test/internal/models"
	"github.com/drag0nfet/effective-mobile-test/internal/repository"
	"github.com/sirupsen/logrus"
)

type PersonService struct {
	repo *repository.PersonRepository
}

func NewPersonService(repo *repository.PersonRepository) *PersonService {
	return &PersonService{repo: repo}
}

func (s *PersonService) CreatePerson(person *models.Person) error {
	logrus.Debugf("Creating person: %+v", person)

	// Здесь будет логика обогащения с помощью API

	err := s.repo.CreatePerson(person)
	if err != nil {
		logrus.Errorf("Failed to create person: %v", err)
		return err
	}
	logrus.Infof("Person created with ID: %d", person.ID)
	return nil
}

func (s *PersonService) GetPeople(filters map[string]interface{}, offset, limit int) ([]models.Person, error) {
	logrus.Debugf("Fetching people with filters: %+v, offset: %d, limit: %d", filters, offset, limit)
	people, err := s.repo.FindPeople(filters, offset, limit)
	if err != nil {
		logrus.Errorf("Failed to fetch people: %v", err)
		return nil, err
	}
	logrus.Infof("Fetched %d people", len(people))
	return people, nil
}

func (s *PersonService) UpdatePerson(person *models.Person) error {
	logrus.Debugf("Updating person with ID: %d", person.ID)
	err := s.repo.UpdatePerson(person)
	if err != nil {
		logrus.Errorf("Failed to update person: %v", err)
		return err
	}
	logrus.Infof("Person with ID %d updated", person.ID)
	return nil
}

func (s *PersonService) DeletePerson(id uint) error {
	logrus.Debugf("Deleting person with ID: %d", id)
	err := s.repo.DeletePerson(id)
	if err != nil {
		logrus.Errorf("Failed to delete person: %v", err)
		return err
	}
	logrus.Infof("Person with ID %d deleted", id)
	return nil
}
