package services

import (
	"bufio"
	"dover/models"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type IngestionService interface {
	IngestBulk(records *bufio.Scanner) error
	Ingest(record string) error
}

type ingestionService struct {
	DB *sqlx.DB
}

func New(d *sqlx.DB) IngestionService {
	return &ingestionService{
		DB: d,
	}
}

func (s *ingestionService) IngestBulk(scanner *bufio.Scanner) error {
	for scanner.Scan() {
		// best effort to ingest all records
		err := s.Ingest(scanner.Text())
		if err != nil {
			// error logged to console and returned later as aggregated
			fmt.Println(err)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (s *ingestionService) Ingest(record string) error {
	_, positions, _, _ := s.parseRecord(record)
	for _, position := range positions {
		fmt.Println(position)
	}
	return nil
}

func (s *ingestionService) parseRecord(raw string) (*models.Profile, []*models.Education, []*models.Position, error) {
	ser := map[string]interface{}{}
	err := json.Unmarshal([]byte(raw), &ser)
	if err != nil {
		return nil, nil, nil, err
	}
	rawProfileJson, ok := ser["profile"].(map[string]interface{})
	if !ok {
		return nil, nil, nil, fmt.Errorf("unable to parse profile to map")
	}
	profile := &models.Profile{
		FirstName:  parseString(rawProfileJson, "firstName"),
		LastName:   parseString(rawProfileJson, "lastName"),
		RawProfile: []byte(raw),
	}
	rawPositionsJson, ok := ser["positions"].([]interface{})
	if !ok {
		return nil, nil, nil, fmt.Errorf("unable to parse positions to list of maps")
	}
	positions := []*models.Position{}
	for _, rp := range rawPositionsJson {
		rawPosition := rp.(map[string]interface{})
		position := &models.Position{
			Title:       parseString(rawPosition, "title"),
			CompanyName: parseString(rawPosition, "companyName"),
		}
		if val, ok := rawPosition["timePeriod"]; ok && val != nil {
			timePeriod := rawPosition["timePeriod"].(map[string]interface{})
			position.StartDate = parseTime(timePeriod, "startDate")
			position.EndDate = parseTime(timePeriod, "endDate")
		}
		positions = append(positions, position)
	}

	rawSkills, ok := ser["skills"].([]interface{})
	if !ok {
		return nil, nil, nil, fmt.Errorf("unable to parse skills to list of maps")
	}
	skills := []string{}
	for _, rs := range rawSkills {
		rawSkill := rs.(map[string]interface{})
		skills = append(skills, parseString(rawSkill, "skillName"))
	}
	profile.Skills = skills
	rawEducations := ser["educations"].([]interface{})
	educations := []*models.Education{}
	for _, re := range rawEducations {
		rawEdu := re.(map[string]interface{})
		education := &models.Education{
			DegreeName:   parseString(rawEdu, "degreeName"),
			FieldOfStudy: parseString(rawEdu, "fieldOfStudy"),
			SchoolName:   parseString(rawEdu, "schoolName"),
		}
		if val, ok := rawEdu["timePeriod"]; ok && val != nil {
			timePeriod := rawEdu["timePeriod"].(map[string]interface{})
			education.StartDate = parseTime(timePeriod, "startDate")
			education.EndDate = parseTime(timePeriod, "endDate")
		}
		educations = append(educations, education)
	}
	return profile, educations, positions, nil
}
