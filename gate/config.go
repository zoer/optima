package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

var _ ConfigServicer = (*ConfigService)(nil)

type configRequestParams struct {
	Type, Data string
}

type ConfigServicer interface {
	HandleParams(http.ResponseWriter, *http.Request)
	SetupRoutes(*mux.Router)
}

// ConfigService represents config params.
type ConfigService struct {
	Repository ConfigRepo // repository where the configs data stored
}

// NewConfigService allocates ConfigService structs.
func NewConfigService(repo ConfigRepo) *ConfigService {
	return &ConfigService{
		Repository: repo,
	}
}

// HandleParams handlers config param requests.
func (c *ConfigService) HandleParams(w http.ResponseWriter, r *http.Request) {
	var data configRequestParams
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, errToJson("Unable to read JSON body"), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if len(data.Type) == 0 || len(data.Data) == 0 {
		http.Error(w, errToJson("'Type' and 'Data' JSON values must be provided"), http.StatusUnprocessableEntity)
		return
	}

	conf, err := c.Repository.GetParams(data.Type, data.Data)
	if err == ErrNotFound {
		http.Error(w, errToJson("Configuration not found"), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, errToJson("Internal error"), http.StatusInternalServerError)
		return
	}

	w.Write(conf)
}

// SetupRoutes sets configs routes.
func (c *ConfigService) SetupRoutes(r *mux.Router) {
	r.Methods("POST").Path("").HandlerFunc(c.HandleParams)
}

// errToJson wraps an error message into JSON structure.
func errToJson(msg string) string {
	out, _ := json.Marshal(map[string]interface{}{"error": msg})
	return string(out)
}
