package businesslayer

import (
	"context"
	"net"

	"github.com/arizon-dread/webdig-backend/models"
)

func LookupDNS(req models.Req) (models.Resp, error) {
	var resp models.Resp

}

func LookupIP(req models.Req) (models.Resp, error) {

}

func getResolver() (net.Resolver, error) {
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d: net.Dialer{
				Timeout: time.Millisecond = time.Duration(5000),f
			}
		},
	}
}
