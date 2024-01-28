package businesslayer

import (
	"context"
	"fmt"
	"net"
	"slices"
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
		err := fmt.Errorf("could not find dns record")
		resp.Err = err
		return resp, err
	} else {
		resp.Err = nil
	}

	for _, ip := range resp.InternalIPAddresses {
		if slices.Contains(resp.ExternalIPAddresses, ip) {
			index := slices.Index(resp.ExternalIPAddresses, ip)
			resp.ExternalIPAddresses = slices.Delete(resp.ExternalIPAddresses, index, index+1)
		}
	}
	resp.DnsNames = slices.Compact(resp.DnsNames)

	return resp, nil
}

func lookupDNS(ctx context.Context, dnsServers []string, isInternal bool, isDNS bool, req models.Req, resp *models.Resp) {
	var wg sync.WaitGroup
	for i, ip := range dnsServers {
		wg.Add(1)
		go func(wg *sync.WaitGroup, i int) {
			defer wg.Done()
			if isDNS {
				ips := lookupIPforDNSandServer(ctx, req.Host, ip, resp)
				if len(ips) > 0 {
					for _, ip := range ips {
						if isInternal {
							if !slices.Contains(resp.InternalIPAddresses, ip.String()) {
								resp.InternalIPAddresses = append(resp.InternalIPAddresses, ip.String())
							}

						} else {
							if !slices.Contains(resp.ExternalIPAddresses, ip.String()) {
								resp.ExternalIPAddresses = append(resp.ExternalIPAddresses, ip.String())
							}
						}
					}
				}
			} else {
				lookupDNSforIpAndServer(ctx, req.Host, ip, resp)

			}
		}(&wg, i)

	}
	wg.Wait()
}

func lookupIPforDNSandServer(ctx context.Context, dnsName string, dnsServer string, resp *models.Resp) []net.IP {

	r, cancel := getResolver(ctx, dnsServer)
	defer cancel()
	ips, err := r.LookupIP(ctx, "ip4", dnsName)
	if err != nil {
		resp.Err = err
		return nil
	} else {
		return ips
	}
}

func lookupDNSforIpAndServer(ctx context.Context, ip string, dnsServer string, resp *models.Resp) {

	r, cancel := getResolver(ctx, dnsServer)
	defer cancel()

	hosts, err := r.LookupAddr(ctx, ip)
	if err != nil {
		resp.Err = err
	} else {
		resp.DnsNames = append(resp.DnsNames, hosts...)
	}

}

func getResolver(ctx context.Context, dnsHost string) (*net.Resolver, context.CancelFunc) {
	ctxTo, cancel := context.WithTimeout(ctx, time.Duration(time.Millisecond*5000))
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Duration(time.Millisecond * 5000),
			}
			return d.DialContext(ctxTo, network, dnsHost+":53")
		},
	}
	return r, cancel
}
