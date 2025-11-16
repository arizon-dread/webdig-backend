package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/arizon-dread/webdig-backend/internal/cliconfig/handlers"
	"github.com/arizon-dread/webdig-backend/pkg/types"
)

func main() {
	var s string
	var c bool
	flag.StringVar(&s, "s", "", "server url to webdig server for this call")
	flag.BoolVar(&c, "c", false, "configure the --server setting in your local config for subsequent usage")
	flag.Parse()
	// get all remaining arguments
	args := flag.Args()
	var addr string
	// make sure that a lookup parameter is supplied
	if len(args) > 1 {
		fmt.Printf("You must supply the lookup address as the last parameter")
		os.Exit(1)
	}
	for _, a := range args {
		if len(a) > 0 {
			addr = a
			break
		}
	}
	if len(addr) == 0 {
		fmt.Println("supply a single IP or DNS FQDN to perform a lookup")
		os.Exit(1)
	}
	conf := &types.ServerConf{}
	// if a server address is supplied, make the request to that server.
	if s != "" {
		// call webdig with specified address
		conf.Server = s
		if c {
			// try to create the config file
			err := handlers.SaveConf(s)
			if err != nil {
				fmt.Printf("%v", err)
				os.Exit(2)
			}
		}
	} else {
		savedConf, err := handlers.EnsureConfig(nil)
		conf.Server = savedConf.Server
		if err != nil {
			fmt.Printf("%v", errors.Unwrap(err))
			os.Exit(3)
		}
	}
	req := types.Req{
		Host: addr,
	}
	resp, err := handlers.MakeCall(req, conf)
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(4)
	}
	var respText string

	for _, r := range resp.Results {
		if len(r.IPAddresses) == 0 && len(r.DnsNames) == 0 {
			continue
		}
		respText += fmt.Sprintf("Name: %v\n", r.Name)
		if len(r.DnsNames) > 0 {
			dnsNames := strings.Join(r.DnsNames, ", ")
			respText += fmt.Sprintf("DNS: %v", dnsNames)
		}
		if len(r.IPAddresses) > 0 {
			ipAddrs := strings.Join(r.IPAddresses, ", ")
			respText += fmt.Sprintf("IP: %v", ipAddrs)
		}
		respText += "\n-----\n"
	}
	fmt.Printf("%v", respText)
	os.Exit(0)
}
