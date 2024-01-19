package businesslayer

import (
	"context"
	"net"
	"time"

	"github.com/arizon-dread/webdig-backend/config"
	"github.com/arizon-dread/webdig-backend/models"
)

func LookupDNS(req models.Req) (models.Resp, error) {
	var resp models.Resp
	cfg := config.GetInstance()
	go func() {
		for _, ip := range cfg.DNS.InternalServers {
			ip, err := lookupDNSforIpAndServer(req.body, ip)
		}
	}
}

func LookupIP(req models.Req) (models.Resp, error) {

}

func lookupDNSforIpAndServer(string ip, string dnsServer) (string, error) {
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
