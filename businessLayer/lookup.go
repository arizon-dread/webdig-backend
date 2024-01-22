package businesslayer

import (
	"context"
	"fmt"
	"net"
	"strings"
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
			if isDNS {
				fmt.Printf("got dns: %v\n", req.Host)
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
						fmt.Printf("ending waitGroup\n")
						wg.Done()
					}
				} else if !isInternal && len(resp.ExternalIPAddresses) > 0 {
					for len(dnsServers)-(i+1) > 0 {
						fmt.Printf("ending waitGroup\n")
						wg.Done()
					}
				}
			} else {
				fmt.Printf("got ip: %v\n", req.Host)
				ptrAddr, err := reverseIPAddress(req.Host)
				fmt.Printf("reversed: %v\n", ptrAddr)
				if err != nil {
					resp.Err = fmt.Errorf("failed to read input as ip address.")
				} else {
					lookupDNSforIpAndServer(ctx, ptrAddr, ip, resp)
					if len(resp.DnsNames) > 0 {
						for len(dnsServers)-(i+1) > 0 {

							fmt.Printf("ending waitGroup\n")
							wg.Done()
						}
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

	hosts, err := r.LookupAddr(ctx, ip)
	if err != nil {
		resp.Err = err
	} else {
		resp.DnsNames = append(resp.DnsNames, hosts...)
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
func reverseIPAddress(ip string) (string, error) {
	var netip = net.ParseIP(ip).To4()
	if netip != nil {
		// split into slice by dot .
		addressSlice := strings.Split(netip.String(), ".")
		reverseSlice := []string{}

		for i := range addressSlice {
			octet := addressSlice[len(addressSlice)-1-i]
			reverseSlice = append(reverseSlice, octet)
		}

		// sanity check
		//fmt.Println(reverseSlice)

		return strings.Join(reverseSlice, "."), nil

	} else {
		return "", fmt.Errorf("invalid ipv4 address")
	}
}
