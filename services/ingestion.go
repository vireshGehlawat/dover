package services

import (
	"bufio"
	"dover/models"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"strings"
)

type IngestionService interface {
	IngestBulk(records *bufio.Scanner) error
	Ingest(record string) error
}

type ingestionService struct {
	DB *sqlx.DB
}

func New(db *sqlx.DB) IngestionService {
	_, err := db.Query("select * from Profile")
	if err != nil {
		if strings.Contains(err.Error(), "doesn't exist") {
			db.MustExec(models.ProfileSchema)
			db.MustExec(models.EducationSchema)
			db.MustExec(models.PositionSchema)
		} else {
			log.Fatal(err)
		}
	}
	return &ingestionService{
		DB: db,
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
	profile, educations, positions, err := s.parseRecord(record)
	// parsing err
	if err != nil {
		return err
	}
	tx := s.DB.MustBegin()
	dbRecord, err := tx.Exec(fmt.Sprintf("INSERT INTO Profile (LastName, FirstName, Skills) VALUES (\"%s\", \"%s\", \"%s\")", profile.LastName, profile.FirstName, profile.Skills))
	if err != nil {
		tx.Rollback()
		fmt.Println(err)
		return err
	}
	profileID, err := dbRecord.LastInsertId()
	for _, e := range educations {
		e.ProfileID = uint64(profileID)
		// todo fix NULLable date fields
		_, err := tx.Exec(fmt.Sprintf(
			"INSERT INTO Education "+
				"(ProfileID, EndDate, StartDate, DegreeName, FieldOfStudy, SchoolName) "+
				"VALUES (\"%d\", \"%s\", \"%s\", \"%s\", \"%s\", \"%s\")",
			e.ProfileID, e.EndDate, e.StartDate, e.DegreeName, e.FieldOfStudy, e.SchoolName),
		)
		if err != nil {
			tx.Rollback()
			fmt.Println(err)
			return err
		}
	}
	for _, p := range positions {
		p.ProfileID = uint64(profileID)
		// todo fix NULLable date fields
		_, err := tx.Exec(fmt.Sprintf(
			"INSERT INTO Position "+
				"(ProfileID, EndDate, StartDate, Title, CompanyName) "+
				"VALUES (\"%d\", \"%s\", \"%s\", \"%s\", \"%s\")",
			p.ProfileID, p.EndDate, p.StartDate, p.Title, p.CompanyName),
		)
		if err != nil {
			tx.Rollback()
			fmt.Println(err)
			return err
		}
	}
	return tx.Commit()
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
	profile.Skills = strings.Join(skills, "__")
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
