package models

type Resp struct {
	Results []Result `json:"result"`
	Err     error    `json:"error"`
}

type Result struct {
	Name        string   `json:"name"`
	DnsNames    []string `json:"dnsNames"`
	IPAddresses []string `json:"ipAddresses"`
	Err         error    `json:"error"`
}
