package types

type Resp struct {
	Results []Result `json:"results"`
	Err     error    `json:"error"`
}

type Result struct {
	Name        string   `json:"name"`
	DnsNames    []string `json:"dnsNames"`
	IPAddresses []string `json:"ipAddresses"`
	Err         error    `json:"error"`
}
