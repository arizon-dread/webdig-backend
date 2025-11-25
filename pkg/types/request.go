package types

type Req struct {
	Host  string `json:"host"`
	CNAME *bool  `json:"cname"`
}
