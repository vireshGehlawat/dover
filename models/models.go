package models

import (
	"database/sql"
)

type AutoIncr struct {
	ID      uint64
}

type Profile struct {
	AutoIncr
	FirstName  string
	LastName   string
	RawProfile sql.RawBytes
	Skills     string
}

type Position struct {
	AutoIncr
	ProfileID    uint64
	Title       string
	CompanyName string
	StartDate   sql.NullTime
	EndDate     sql.NullTime
}

type Education struct {
	AutoIncr
	ProfileID    uint64
	DegreeName   string
	StartDate    sql.NullTime
	EndDate      sql.NullTime
	FieldOfStudy string
	SchoolName   string
}
