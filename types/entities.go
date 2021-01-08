package types

import "github.com/pborman/uuid"

type ListViewFilters struct {
	HasCSDegree *bool
	IsEmployed  *bool
	Offset      int32
	Limit       int32
}

//ProfileListView is the wrapper for the list view
type ProfileListView struct {
	ViewID          uuid.UUID
	FullName        string
	CurrentTitle    string
	TotalExperience int32
	HasCSDegree     *bool
	IsEmployed      *bool
}
