package businesslayer

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/arizon-dread/webdig-backend/config"
	"github.com/arizon-dread/webdig-backend/models"
)

func LookupDNS(req models.Req) (models.Resp, error) {
	var resp models.Resp
	cfg := config.GetInstance()
	var wg sync.WaitGroup
	wg.Add(len(cfg.DNS.InternalServers))
	for i, ip := range cfg.DNS.InternalServers {
		go func(wg *sync.WaitGroup, i int) {
			dns, err := lookupDNSforIpAndServer(req.Host, ip)
			if err == nil {
				resp.DnsNames = append(resp.DnsNames, dns)
				for len(cfg.DNS.InternalServers)-(i+1) > 0 {
					wg.Done()
				}
				return
			}
			wg.Done()
		}(&wg, i)
	}
	wg.Wait()
	if len(resp.DnsNames) == 0 {
		err := fmt.Errorf("could not find internal dns record")
		resp.Err = err
		return resp, err
	}
	return resp, nil
}

func LookupIP(req models.Req) (models.Resp, error) {

}

func lookupDNSforIpAndServer(ip string, dnsServer string) (string, error) {
	r := getResolver(dnsServer)

}

func getResolver(dnsHost string) *net.Resolver {
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
