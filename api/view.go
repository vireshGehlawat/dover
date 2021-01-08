package api

import (
	"dover/services"
	"dover/types"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type View interface {
	GetProfileList(writer http.ResponseWriter, request *http.Request)
	GetProfile(writer http.ResponseWriter, request *http.Request)
}

type view struct {
	profiles services.ProfilesService
}

func New(profiles services.ProfilesService) View {
	return &view{
		profiles: profiles,
	}
}

func (v *view) GetProfileList(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("received GetProfileList req")
	err := request.ParseForm()
	if err != nil {
		// error in parsing the request
		fmt.Println(err)
		writer.WriteHeader(400)
		return
	}
	filters := v.parseFiltersFromRequest(request)
	fmt.Println(filters)
	// todo get this from DB after processing filters
	profiles, err := v.profiles.GetProfiles(request.Context(), v.parseFiltersFromRequest(request))
	if err != nil {
		fmt.Println(err)
		writer.WriteHeader(500)
		return
	}
	v.respondWithJSON(writer, 200, profiles)
}

func (v *view) parseFiltersFromRequest(*http.Request) *types.ListViewFilters {
	return &types.ListViewFilters{}
}

func (v *view) GetProfile(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("received GetProfile req")
	err := request.ParseForm()
	if err != nil {
		fmt.Println("error in parsing the request")
		writer.WriteHeader(400)
		return
	}
	id := request.Form.Get("id")
	profileID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		fmt.Println(err)
		writer.WriteHeader(400)
		return
	}
	// todo get this from DB after processing
	profile, err := v.profiles.GetProfile(request.Context(), profileID)
	if err != nil {
		fmt.Println(err)
		writer.WriteHeader(500)
		return
	}
	v.respondWithJSON(writer, 200, profile)
}

func (v *view) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
