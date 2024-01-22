package businesslayer

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
	"unicode"

	"github.com/arizon-dread/webdig-backend/config"
	"github.com/arizon-dread/webdig-backend/models"
)

func Lookup(ctx context.Context, req models.Req) (models.Resp, error) {
	var resp models.Resp
	cfg := config.GetInstance()
	isDNS := func() bool {
		for _, r := range req.Host {
			if unicode.IsLetter(r) {
				return true
			}
		}
		return false
	}

	lookupDNS(ctx, cfg.DNS.InternalServers, true, isDNS(), req, &resp)

	lookupDNS(ctx, cfg.DNS.ExternalServers, false, isDNS(), req, &resp)

	if len(resp.DnsNames) == 0 && len(resp.ExternalIPAddresses) == 0 && len(resp.InternalIPAddresses) == 0 {
		err := fmt.Errorf("could not find internal dns record")
		resp.Err = err
		return resp, err
	}
	return resp, nil
}

func lookupDNS(ctx context.Context, dnsServers []string, isInternal bool, isDNS bool, req models.Req, resp *models.Resp) {
	var wg sync.WaitGroup
	wg.Add(len(dnsServers))
	for i, ip := range dnsServers {
		go func(wg *sync.WaitGroup, i int) {
			if !isDNS {
				lookupDNSforIpAndServer(ctx, req.Host, ip, resp)
				if len(resp.DnsNames) > 0 {
					for len(dnsServers)-(i+1) > 0 {
						wg.Done()
					}
				}
			} else {
				ips := LookupIPforDNSandServer(ctx, req.Host, ip, resp)
				if len(ips) > 0 {
					for _, ip := range ips {
						if isInternal {
							resp.InternalIPAddresses = append(resp.InternalIPAddresses, ip.String())
						} else {
							resp.ExternalIPAddresses = append(resp.ExternalIPAddresses, ip.String())
						}
					}
				}
				if isInternal && len(resp.InternalIPAddresses) > 0 {
					for len(dnsServers)-(i+1) > 0 {
						wg.Done()
					}
				} else if !isInternal && len(resp.ExternalIPAddresses) > 0 {
					for len(dnsServers)-(i+1) > 0 {
						wg.Done()
					}
				}
			}

			wg.Done()
		}(&wg, i)

	}
	wg.Wait()
}

func LookupIPforDNSandServer(ctx context.Context, dnsName string, dnsServer string, resp *models.Resp) []net.IP {

	r := getResolver(ctx, dnsServer)
	ips, err := r.LookupIP(ctx, "ip4", dnsName)
	if err != nil {
		resp.Err = err
		return nil
	} else {
		return ips
	}
}

func lookupDNSforIpAndServer(ctx context.Context, ip string, dnsServer string, resp *models.Resp) {

	r := getResolver(ctx, dnsServer)

	host, err := r.LookupCNAME(ctx, ip)
	if err != nil {
		resp.Err = err
	} else {
		resp.DnsNames = append(resp.DnsNames, host)
	}

}

func getResolver(ctx context.Context, dnsHost string) *net.Resolver {
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Duration(time.Millisecond * 5000),
			}
			return d.DialContext(ctx, network, dnsHost+":53")
		},
	}
	return r
}
