package businesslayer

import (
	"context"
	"net"
	"time"

	"github.com/arizon-dread/webdig-backend/models"
)

func LookupDNS(req models.Req) (models.Resp, error) {
	var resp models.Resp

}

func LookupIP(req models.Req) (models.Resp, error) {

}

func lookupDNSforIpAndServer() (string, error) {

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
