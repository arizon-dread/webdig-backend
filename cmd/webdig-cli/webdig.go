package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/arizon-dread/webdig-backend/internal/cliconfig/handlers"
	"github.com/arizon-dread/webdig-backend/pkg/types"
)

var (
	s = flag.String("s", "", "URL to webdig SERVER for this call, including contextpath. Ex: -s http://localhost:8080/api/dig")
	c = flag.Bool("c", false, "Configure the server setting in your local config for subsequent requests, -s http://webdig.server.name/api/dig is required to use -c.\nThe server.yaml config file is stored in:\nLinux: /home/<username>/.config/webdig/\nMacOS: /Users/<username>/Library/Application Support/webdig/\nWindows: C:\\Users\\<username>\\AppData\\webdig\\\n")
	n = flag.Bool("n", false, "Lookup cName. Using this flag makes the command slower but also looks up the CNAME (if any) and which A-record it points to\n")
)

func main() {
	flag.Usage = setHelpText
	flag.Parse()
	// get all remaining arguments
	args := flag.Args()
	var addr string
	// make sure that a lookup parameter is supplied
	if len(args) > 1 {
		fmt.Println("You must supply the lookup address as the last parameter")
		os.Exit(1)
	}
	for _, a := range args {
		if len(a) > 0 {
			addr = a
			break
		}
	}

	if len(addr) == 0 {
		fmt.Println("Supply a single IP or DNS FQDN to perform a lookup")
		os.Exit(1)
	}
	conf := &types.ServerConf{}

	// if a server address is supplied, make the request to that server.
	if *s != "" {
		// call webdig with specified address
		conf.Server = *s
		if *c {
			// try to create the config file
			err := handlers.SaveConf(*s)
			if err != nil {
				fmt.Printf("Error trying to save your config, %v\n", err)
				os.Exit(2)
			}
		}
	} else {
		savedConf, err := handlers.EnsureConfig(nil)
		if err != nil {
			fmt.Printf("Couldn't find configured service, use -c -s full-url-to-server to configure, see -h for help %v\n", err)
			os.Exit(3)
		}
		conf.Server = savedConf.Server
	}
	req := types.Req{
		Host:  addr,
		CNAME: n,
	}
	resp, err := handlers.MakeCall(req, conf)
	if err != nil {
		// Error formatting is done inside handlers.MakeCall(), just output whatever text we get in the err.
		fmt.Printf("%v", err)
		os.Exit(4)
	}
	var respText strings.Builder

	for _, r := range resp.Results {
		if len(r.IPAddresses) == 0 && len(r.DnsNames) == 0 {
			continue
		}
		fmt.Fprintf(&respText, "Name: %v\n", r.Name)
		if len(r.DnsNames) > 0 {
			dnsNames := strings.Join(r.DnsNames, ", ")
			fmt.Fprintf(&respText, "DNS: %v", dnsNames)
		}
		if len(r.IPAddresses) > 0 {
			ipAddrs := strings.Join(r.IPAddresses, ", ")
			fmt.Fprintf(&respText, "IP: %v", ipAddrs)
		}
		if len(r.Cname) > 0 {
			fmt.Fprintf(&respText, "\nCNAME for %v", r.Cname)
		}
		respText.WriteString("\n-----\n")
	}
	fmt.Printf("\n%v", respText.String())
	os.Exit(0)
}

func setHelpText() {
	fmt.Fprintf(flag.CommandLine.Output(), `wdig - a commandline client for the WebDig API

Usage: %s [options] <arg>

options:

`, os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintf(flag.CommandLine.Output(), `
Lookup: 
	<arg> is the ip or dnsName you want to perform lookup on. 

You must point your wdig client to a WebDig API, example:

	wdig -c -s http://localhost:8080/api/dig

You can also skip -c to just run wdig towards a different API address on this call, without persisting the address to your config.
After configuring your client, perform lookups like this:

	wdig www.google.com
	wdig 192.168.1.1
	wdig -n www.google.com
	wdig -s http://10.10.33.33/api/dig www.google.com`)
	os.Exit(0)
}
