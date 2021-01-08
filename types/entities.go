package types

type ListViewFilters struct {
	HasCSDegree *bool
	IsEmployed  *bool
	Offset      int32
	Limit       int32
}

//ProfileListView is the wrapper for the list view
type ProfileListView struct {
	ProfileID       int64
	FullName        string
	CurrentTitle    string
	TotalExperience int32
	HasCSDegree     *bool
	IsEmployed      *bool
}
