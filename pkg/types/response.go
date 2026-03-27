package types

type Resp struct {
	Results []Result `json:"results"`
	Err     error    `json:"error"`
}

type Result struct {
	Name        string      `json:"name"`
	DnsNames    []string    `json:"dnsNames"`
	IPAddresses []string    `json:"ipAddresses"`
	Cname       string      `json:"cnameFor"`
	Type        string      `json:"type"`
	Err         interface{} `json:"error"`
}
