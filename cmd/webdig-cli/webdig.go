package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/arizon-dread/webdig-backend/internal/cliconfig/handlers"
	"github.com/arizon-dread/webdig-backend/pkg/types"
)

func main() {
	var s string
	var c bool
	flag.StringVar(&s, "s", "", "server url to webdig server for this call")
	args := flag.Args()
	flag.BoolVar(&c, "c", false, "configure the --server setting in your local config for subsequent usage")
	flag.Parse()

	// make sure that a lookup parameter is supplied
	if len(args) != 1 {
		fmt.Println("supply a single IP or DNS FQDN to perform a lookup")
		os.Exit(1)
	}
	// if a server address is supplied, make the request to that server.
	if s != "" {
		// call webdig with specified address
		if c {
			// try to create the config file
			err := handlers.SaveConf(s)
			if err != nil {
				fmt.Printf("%v", err)
				os.Exit(2)
			}
		}
	}
	conf, err := handlers.EnsureConfig(nil)
	if err != nil {
		fmt.Printf("%v", errors.Unwrap(err))
		os.Exit(3)
	}

	req := types.Req{
		Host: args[0],
	}
	resp, err := makeCall(req, conf)
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(4)
	}
	var respText string

	for _, r := range resp.Results {
		respText += fmt.Sprintf("Name: %v ", r.Name)
		if len(r.DnsNames) > 0 {
			dnsNames := strings.Join(r.DnsNames, ", ")
			respText += fmt.Sprintf("DNS: %v", dnsNames)
		}
		if len(r.IPAddresses) > 0 {
			ipAddrs := strings.Join(r.IPAddresses, ", ")
			respText += fmt.Sprintf("IP: %v", ipAddrs)
		}
	}
}

func makeCall(req types.Req, conf *types.ServerConf) (*types.Resp, error) {
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
	var resp *types.Resp
	var b []byte
	_, err = res.Body.Read(b)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body, %w", err)
	}
	err = json.Unmarshal(b, resp)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal response into go struct, %w", err)
	}
	return resp, nil
}
