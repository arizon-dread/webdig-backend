package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/arizon-dread/webdig-backend/internal/cliconfig/platform"
	"github.com/arizon-dread/webdig-backend/internal/dns"
	"github.com/arizon-dread/webdig-backend/pkg/types"
	"gopkg.in/yaml.v3"
)

func main() {
	var s string
	var c *bool
	flag.StringVar(&s, "server", "", "url to webdig server for this call")
	args := flag.Args()
	flag.BoolVar(c, "config", false, "configure the --server setting in your local config for subsequent usage")
	flag.Parse()
	if s != "" {
		pFind := platform.NewFindPath()
		confPath := pFind.FindPath()
		fullConfPath := fmt.Sprintf("%v%vwebdig%vserver.yaml", confPath, os.PathSeparator, os.PathSeparator)
		err := os.Mkdir(fullConfPath, 0o755)
		if err != nil {
			fmt.Printf("unable to create config dir, %v", err)
		}
		f, err := os.Open(fullConfPath)
		if err != nil {
			fmt.Printf("unable to open file, %v", fullConfPath)
			os.Exit(1)
		}
		defer f.Close()

		var b []byte
		_, err = f.Read(b)
		if err != nil {
			fmt.Printf("unable to read file, %v", err)
			os.Exit(2)
		}
		if s != "" {
			err = saveConf(s, f)
			if err != nil {
				fmt.Printf("%v", err)
				os.Exit(3)
			}
		}

	}

	if len(args) != 1 {
		fmt.Println("supply a single IP or DNS FQDN to perform a lookup")
		os.Exit(4)
	}
	req := types.Req{
		Host: args[0],
	}
	resp, err := dns.Lookup(context.Background(), req)
	if err != nil {
		fmt.Printf("error during lookup: %v\n", err)
		if resp.Err != nil {
			fmt.Printf("errors from lookup: %v\n", resp.Err)
		}
		os.Exit(5)
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

func saveConf(url string, f *os.File) error {
	var (
		ErrUnmarshal = errors.New("unmarshal error")
		ErrSaveConf  = errors.New("error saving config")
	)
	conf := types.ServerConf{
		Server: url,
	}
	c, err := yaml.Marshal(conf)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrUnmarshal, err)
	}
	_, err = f.Write(c)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSaveConf, err)
	}
	return nil
}
