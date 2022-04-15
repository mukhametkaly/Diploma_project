package auth

import "github.com/mukhametkaly/Diploma_project/auth-api/src/models"

var (
	AdminRights = models.Rights{
		Products: 0,
		Category: 0,
		Merchant: 0,
	}

	CashierRights = models.Rights{
		Products: 0,
		Category: 0,
		Merchant: 0,
	}

	MerchantRights = models.Rights{
		Products: 0,
		Category: 0,
		Merchant: 0,
	}
)
