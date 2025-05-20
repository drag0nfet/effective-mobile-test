package service

import (
	"encoding/json"
	"errors"
	"github.com/drag0nfet/effective-mobile-test/internal/config"
	"github.com/drag0nfet/effective-mobile-test/internal/models"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type EnrichmentService struct {
	cfg *config.Config
}

func NewEnrichmentService(cfg *config.Config) *EnrichmentService {
	return &EnrichmentService{cfg: cfg}
}

func (e *EnrichmentService) EnrichPerson(person *models.Person) error {
	logrus.Debugf("EnrichPerson called")
	// Запрос возраста (Agify API)
	err := e.getAge(person)
	if err != nil {
		return err
	}

	// Запрос пола (Genderize API)
	err = e.getGender(person)
	if err != nil {
		return err
	}

	// Запрос национальности (Nationalize API)
	err = e.getNationalize(person)
	if err != nil {
		return err
	}

	logrus.Infof("Person with ID %d enriched successfully", person.ID)
	return nil
}

func (e *EnrichmentService) getAge(person *models.Person) error {
	ageResp, err := http.Get(e.cfg.AgifyAPIURL + "?name=" + person.Name)
	if err != nil {
		logrus.Errorf("Failed to fetch age from Agify: %v", err)
		return err
	}

	defer ageResp.Body.Close()

	ageBody, err := io.ReadAll(ageResp.Body)
	if err != nil {
		logrus.Errorf("Failed to read age response: %v", err)
		return err
	}

	var ageData struct {
		Age int `json:"age"`
	}

	if err = json.Unmarshal(ageBody, &ageData); err == nil && ageData.Age > 0 {
		person.Age = ageData.Age
	} else if err != nil {
		logrus.Errorf("Failed to unmarshall age: %v", err)
		return err
	} else {
		logrus.Errorf("Incorrect age got from API: %d", ageData.Age)
		return errors.New("incorrect age from API")
	}

	logrus.Debugf("Age found correctly")
	return nil
}

func (e *EnrichmentService) getGender(person *models.Person) error {
	genderResp, err := http.Get(e.cfg.GenderizeAPIURL + "?name=" + person.Name)
	if err != nil {
		logrus.Errorf("Failed to fetch gender from Genderize: %v", err)
		return err
	}

	defer genderResp.Body.Close()

	genderBody, err := io.ReadAll(genderResp.Body)
	if err != nil {
		logrus.Errorf("Failed to read gender response: %v", err)
		return err
	}

	var genderData struct {
		Gender string `json:"gender"`
	}

	if err = json.Unmarshal(genderBody, &genderData); err == nil && genderData.Gender != "" {
		person.Gender = genderData.Gender
	} else if err != nil {
		logrus.Errorf("Failed to unmarshall gender: %v", err)
		return err
	} else {
		logrus.Error("Empty gender got from API")
		return errors.New("empty gender from API")
	}

	logrus.Debugf("Gender found correctly")
	return nil
}

func (e *EnrichmentService) getNationalize(person *models.Person) error {
	nationalityResp, err := http.Get(e.cfg.NationalizeAPIURL + "?name=" + person.Name)
	if err != nil {
		logrus.Errorf("Failed to fetch nationality from Nationalize: %v", err)
		return err
	}

	defer nationalityResp.Body.Close()

	nationalityBody, err := io.ReadAll(nationalityResp.Body)
	if err != nil {
		logrus.Errorf("Failed to read nationality response: %v", err)
		return err
	}

	var nationalityData struct {
		Country []struct {
			CountryID   string  `json:"country_id"`
			Probability float64 `json:"probability"`
		} `json:"country"`
	}

	if err = json.Unmarshal(nationalityBody, &nationalityData); err == nil && len(nationalityData.Country) > 0 {
		person.Nationality = nationalityData.Country[0].CountryID
	} else if err != nil {
		logrus.Errorf("Failed to unmarshall nationality: %v", err)
		return err
	} else {
		logrus.Errorf("Too many countries got from API: %d pcs", len(nationalityData.Country))
		return errors.New("too many countries from API")
	}

	logrus.Debugf("Nationality found correctly")
	return nil
}
