package models

import (
	"database/sql"
)

type Profile struct {
	ID         uint64       `db:"ID"`
	FirstName  string       `db:"FirstName"`
	LastName   string       `db:"LastName"`
	RawProfile sql.RawBytes `db:"RawProfile"`
	Skills     string       `db:"Skills"`
}

type Position struct {
	ID          uint64 `db:"ID"`
	ProfileID   uint64 `db:"ProfileID"`
	Title       string `db:"Title"`
	CompanyName string `db:"CompanyName"`
	// todo use time type
	StartDate string `db:"StartDate"`
	EndDate   string `db:"EndDate"`
}

type Education struct {
	ID         uint64 `db:"ID"`
	ProfileID  uint64 `db:"ProfileID"`
	DegreeName string `db:"DegreeName"`
	// todo use time type
	StartDate    string `db:"StartDate"`
	EndDate      string `db:"EndDate"`
	FieldOfStudy string `db:"FieldOfStudy"`
	SchoolName   string `db:"SchoolName"`
}
