package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/arizon-dread/webdig-backend/internal"
	"github.com/arizon-dread/webdig-backend/pkg/types"
)

func Lookup(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not read request body"))
		return
	}
	defer r.Body.Close()
	var req types.Req
	json.Unmarshal(b, &r)

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
func Healthz(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Healthy"))
}
