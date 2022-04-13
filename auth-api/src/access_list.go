package src

import "github.com/mukhametkaly/Diploma/auth-api/models"

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
