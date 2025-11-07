package internal

import (
	"context"
	"fmt"
	"net"
	"slices"
	"sync"
	"time"
	"unicode"

	"github.com/arizon-dread/webdig-backend/config"
	"github.com/arizon-dread/webdig-backend/pkg/types"
)

// Evaluate if req.Host is an IP address or a DNS record, lookup the matching type by calling the implementation
// specific to that type of request.
func Lookup(ctx context.Context, req types.Req) (types.Resp, error) {
	var resp types.Resp
	cfg := config.GetInstance()
	isDNS := func() bool {
		for _, r := range req.Host {
			if unicode.IsLetter(r) {
				return true
			}
		}
		return false
	}
	for _, serverGroup := range cfg.DNS {
		lookupDNS(ctx, serverGroup, isDNS(), req, &resp)
	}

	//Find out if we have a no-hit result and return 404.
	none := true
	if isDNS() {
		for _, result := range resp.Results {
			if len(result.IPAddresses) > 0 {
				none = false
				break
			}
		}
	} else {
		for _, result := range resp.Results {
			if len(result.DnsNames) > 0 {
				none = false
				break
			}
		}
	}
	if none {
		err := fmt.Errorf("could not find dns record")
		resp.Err = err
		return resp, err
	}

	removeDuplicates(isDNS(), &resp)

	return resp, nil
}

func removeDuplicates(isDNS bool, resp *types.Resp) {
	//remove duplicates if configured
	cfg := config.GetInstance()
	type dnsContent struct {
		addresses []string
		filter    []string
	}
	dns := make(map[string]dnsContent)
	for _, grp := range cfg.DNS {
		cnt := dnsContent{addresses: grp.Servers, filter: grp.FilterDuplicates}
		dns[grp.Name] = cnt
	}
	rslt := make(map[string][]string)
	if isDNS {
		for _, res := range resp.Results {
			slices.Sort(res.IPAddresses)
			res.IPAddresses = slices.Compact(res.IPAddresses)
			rslt[res.Name] = res.IPAddresses
		}
	} else {

		for _, res := range resp.Results {
			slices.Sort(res.DnsNames)
			res.DnsNames = slices.Compact(res.DnsNames)

			rslt[res.Name] = res.DnsNames
		}
	}
	for k, v := range dns {
		if len(v.filter) > 0 {
			rslt[k] = slices.DeleteFunc[[]string, string](rslt[k], func(e string) bool {
				for _, res := range resp.Results {
					if res.Name != k && slices.Contains(dns[k].filter, res.Name) { //don't check duplicates on the current node k within the current result, then it'd be self comparing, check if the current node has a filter entry with the name of the res
						if isDNS {
							return slices.Contains(res.IPAddresses, e)
						} else {
							return slices.Contains(res.DnsNames, e)
						}
					}
				}
				return false
			})
		}
	}
	for i, v := range resp.Results {
		if isDNS {
			resp.Results[i].IPAddresses = rslt[v.Name]
		} else {
			resp.Results[i].DnsNames = rslt[v.Name]
		}
	}
}

func lookupDNS(ctx context.Context, serverGroup config.ServerGroup, isDNS bool, req types.Req, resp *types.Resp) {
	result := types.Result{
		Name: serverGroup.Name,
	}
	var wg sync.WaitGroup
	for i, ip := range serverGroup.Servers {
		wg.Add(1)
		go func(wg *sync.WaitGroup, i int, ip string) {
			defer wg.Done()
			if isDNS {
				ips, err := lookupIPforDNSandServer(ctx, req.Host, ip)
				if err != nil {
					result.Err = err
				}
				if len(ips) > 0 {
					for _, ip := range ips {
						if !slices.Contains(result.IPAddresses, ip.String()) {
							result.IPAddresses = append(result.IPAddresses, ip.String())
						}
					}
				}
			} else {
				hosts, err := lookupDNSforIpAndServer(ctx, req.Host, ip)
				if err != nil {
					result.Err = err
				}
				if len(hosts) > 0 {
					result.DnsNames = append(result.DnsNames, hosts...)
				}

			}
		}(&wg, i, ip)

	}
	wg.Wait()
	resp.Results = append(resp.Results, result)
}

func lookupIPforDNSandServer(ctx context.Context, dnsName string, dnsServer string) ([]net.IP, error) {
	r, cancel := getResolver(ctx, dnsServer)
	defer cancel()
	ips, err := r.LookupIP(ctx, "ip4", dnsName)
	if err != nil {
		if len(ips) > 0 {
			return ips, err
		}
		return nil, err
	}
	return ips, nil
}

func lookupDNSforIpAndServer(ctx context.Context, ip string, dnsServer string) ([]string, error) {

	r, cancel := getResolver(ctx, dnsServer)
	defer cancel()

	return r.LookupAddr(ctx, ip)

}

// Gets a net.Dialer inside a net.Resolver to perform dns lookup
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
