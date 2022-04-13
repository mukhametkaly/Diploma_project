package models

type User struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Mobile     string `json:"mobile"`
	Salt       string `json:"salt"`
	Role       string `json:"role"`
	ACL        ACL    `json:"acl"`
	MerchantId string `json:"merchantId"`
}

type ACL struct {
	Rights map[string]Rights
}

type Rights struct {
	Products int `json:"products"`
	Category int `json:"category"`
	Merchant int `json:"merchant"`
}
