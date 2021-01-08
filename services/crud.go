package services

import (
	"context"
	"dover/models"
	"dover/types"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

type ProfilesService interface {
	GetProfiles(ctx context.Context, filters *types.ListViewFilters) ([]*types.ProfileListView, error)
	GetProfile(ctx context.Context, profile int64) (*types.ProfileListView, error)
}

type profilesService struct {
	DB *sqlx.DB
}

func NewProfilesService(DB *sqlx.DB) ProfilesService {
	return &profilesService{
		DB: DB,
	}
}

func (s *profilesService) GetProfiles(ctx context.Context, filters *types.ListViewFilters) ([]*types.ProfileListView, error) {
	profiles, err := s.getProfiles()
	if err != nil {
		return nil, err
	}
	profileList := []*types.ProfileListView{}
	for _, p := range profiles {
		profileListView, err := s.prepareView(p)
		if err != nil {
			return profileList, err
		}
		profileList = append(profileList, profileListView)
	}
	return profileList, nil
}

func (s *profilesService) prepareView(p models.Profile) (*types.ProfileListView, error) {
	positions, err := s.getPositions(p.ID)
	if err != nil {
		return nil, err
	}
	educations, err := s.getEducations(p.ID)
	if err != nil {
		return nil, err
	}
	hasCSDegree := s.getHasCSDegree(educations)
	currentTitle := s.getCurrentTitle(positions)
	isEmployed := currentTitle != ""
	profileListView := &types.ProfileListView{
		FullName:        fmt.Sprintf("%s %s", p.FirstName, p.LastName),
		CurrentTitle:    currentTitle,
		TotalExperience: s.getTotalExp(positions),
		HasCSDegree:     &hasCSDegree,
		IsEmployed:      &isEmployed,
		ProfileID:       int64(p.ID),
	}
	return profileListView, nil
}

func (s *profilesService) getProfiles() ([]models.Profile, error) {
	profiles := []models.Profile{}
	rows, err := s.DB.Queryx("SELECT * FROM Profile")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		profile := models.Profile{}
		err := rows.StructScan(&profile)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}
	return profiles, nil
}

func (s *profilesService) getPositions(pid uint64) ([]models.Position, error) {
	positions := []models.Position{}
	rows, err := s.DB.Queryx(fmt.Sprintf("SELECT * FROM Position where ProfileID = %d", pid))
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		p := models.Position{}
		err := rows.StructScan(&p)
		if err != nil {
			return nil, err
		}
		positions = append(positions, p)
	}
	return positions, nil
}

func (s *profilesService) getEducations(pid uint64) ([]models.Education, error) {
	educations := []models.Education{}
	rows, err := s.DB.Queryx(fmt.Sprintf("SELECT * FROM Education where ProfileID = %d", pid))
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		e := models.Education{}
		err := rows.StructScan(&e)
		if err != nil {
			return nil, err
		}
		educations = append(educations, e)
	}
	return educations, nil
}

func (s *profilesService) GetProfile(ctx context.Context, profile int64) (*types.ProfileListView, error) {
	prof := models.Profile{}
	err := s.DB.Get(&prof, fmt.Sprintf("SELECT * FROM Position where ID = %d"))
	if err != nil {
		return nil, err
	}
	return s.prepareView(prof)
}

func (c *profilesService) getCurrentTitle(positions []models.Position) string {
	for _, p := range positions {
		if !validTime(p.EndDate) {
			return p.Title
		}
	}
	return ""
}

func (c *profilesService) getTotalExp(positions []models.Position) int32 {
	totalDays := int32(0)
	for _, p := range positions {
		// if end date is valid use time diff else take time difference from now
		if validTime(p.EndDate) {
			endDateTime, _ := time.Parse("2006-01-02", p.EndDate)
			startDate, _ := time.Parse("2006-01-02", p.StartDate)
			delta := endDateTime.Sub(startDate)
			totalDays += int32(delta.Hours()) / 24
		} else {
			startDate, _ := time.Parse("2006-01-02", p.StartDate)
			delta := time.Now().Sub(startDate)
			totalDays += int32(delta.Hours()) / 24
		}
	}
	return totalDays / 365
}

// todo fix this to be in line with time.Time
func validTime(endDate string) bool {
	return endDate != "0001-01-01"
}

func (c *profilesService) getHasCSDegree(educations []models.Education) bool {
	for _, e := range educations {
		if strings.Contains(strings.ToLower(e.FieldOfStudy), "computer") {
			return true
		}
	}
	return false
}
