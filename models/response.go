package models

type Resp struct {
	DnsNames            []string `json:"dnsNames"`
	ExternalIPAddresses []string `json:"externalIPAddresses"`
	InternalIPAddresses []string `json:"internalIPAddresses"`
	Err                 error    `json:"error"`
}
