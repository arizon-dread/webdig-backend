package models

type Resp struct {
	DnsNames    []string `json: dnsNames`
	IpAddresses []string `json: ipAddresses`
	Err         error    `json: error`
}
