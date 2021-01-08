package api

import (
	"dover/types"
	"encoding/json"
	"fmt"
	"github.com/pborman/uuid"
	"net/http"
)

type View interface {
	GetProfileList(writer http.ResponseWriter, request *http.Request)
	GetProfile(writer http.ResponseWriter, request *http.Request)
}

type view struct {
}

func New() View {
	return &view{}
}

func (v *view) GetProfileList(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("received GetProfileList req")
	err := request.ParseForm()
	if err != nil {
		// error in parsing the request
		writer.WriteHeader(400)
		return
	}
	filters := v.parseFiltersFromRequest(request)
	fmt.Println(filters)
	// todo get this from DB after processing filters
	profiles := []*types.ProfileListView{}
	v.respondWithJSON(writer, 200, profiles)
}

func (v *view) parseFiltersFromRequest(*http.Request) *types.ListViewFilters {
	return &types.ListViewFilters{}
}

func (v *view) GetProfile(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("received GetProfile req")
	err := request.ParseForm()
	if err != nil {
		// error in parsing the request
		writer.WriteHeader(400)
		return
	}
	id := request.Form.Get("id")
	viewID := uuid.Parse(id)
	if viewID.String() == uuid.NIL.String() {
		// error in parsing the request
		writer.WriteHeader(400)
		return
	}
	// todo get this from DB after processing
	profile := &types.ProfileListView{}
	v.respondWithJSON(writer, 200, profile)
}

func (v *view) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
