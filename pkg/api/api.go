package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/arizon-dread/webdig-backend/config"
	"github.com/arizon-dread/webdig-backend/internal"
	"github.com/arizon-dread/webdig-backend/pkg/types"
)

// Unmarshal request, call internal.Lookup to lookup the matching DNS<->IP.
// Return a 200 OK with matches for each server group if the address is found.
// Return a 400 Bad Request if the request cannot be marshalled into a Go struct.
// Return a 404 Not Found if no address is found.
// Return a 500 Internal Server Error if there's an error marshalling the response from the DNS servers into JSON.
func Lookup(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not read request body"))
		return
	}
	defer r.Body.Close()
	var req types.Req
	json.Unmarshal(b, &req)

	var status int = http.StatusOK
	resp, err := internal.Lookup(r.Context(), req)
	if err != nil {
		if strings.Contains(err.Error(), "could not find dns record") {
			status = http.StatusNotFound
		} else {
			status = http.StatusBadRequest
		}

	}
	respJSON, err := json.Marshal(resp)
	if err != nil {
		status = http.StatusInternalServerError
		respJSON = []byte("{\"error\": \"unable to marshal dns query response\"}")
	}
	w.WriteHeader(status)
	w.Write([]byte(respJSON))

}

// Kubernetes healthcheck endpoint.
func Healthz(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Healthy"))
}

func Version(w http.ResponseWriter, r *http.Request) {
	c := config.GetInstance()
	w.Write([]byte(fmt.Sprintf("{\"version\": \"%v\"}", c.General.Version)))
}
