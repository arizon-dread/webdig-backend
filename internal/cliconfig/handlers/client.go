package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/arizon-dread/webdig-backend/pkg/types"
)

func MakeCall(req types.Req, conf *types.ServerConf) (*types.Resp, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling request to json payload, %w", err)
	}
	r, err := http.NewRequest("POST", conf.Server, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("unable to create http request, %w", err)
	}
	r.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return nil, fmt.Errorf("error during lookup: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned statuscode: %d", res.StatusCode)
	}
	resp := &types.Resp{}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body, %w", err)
	}
	err = json.Unmarshal(b, resp)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal response into go struct, %w", err)
	}
	return resp, nil
}
