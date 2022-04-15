package models

type User struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	FullName   string `json:"full_name"`
	IIN        string `json:"IIN"`
	Mail       string `json:"mail"`
	Mobile     string `json:"mobile"`
	Salt       string `json:"salt"`
	Role       string `json:"role"`
	ACL        ACL    `json:"acl"`
	MerchantId string `json:"merchant_id"`
}

type ACL struct {
	Rights map[string]Rights
}

type Rights struct {
	Products int `json:"products"`
	Category int `json:"category"`
	Merchant int `json:"merchant"`
}
