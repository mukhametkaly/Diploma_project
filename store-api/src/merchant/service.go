package merchant

import (
	"context"
	"github.com/mukhametkaly/Diploma/store-api/src/models"
	"math/rand"
	"time"
)

type service struct {
}

// Service is the interface that provides methods.
type Service interface {
	CreateMerchant(request models.Merchant) (models.Merchant, error)
	UpdateMerchant(request models.Merchant) (models.Merchant, error)
	DeleteByIdMerchant(id string) error
	DeleteBatchMerchant(ids []string) error
	GetMerchantById(id string) (models.Merchant, error)
	FilterMerchants(request FilterMerchantsRequest) ([]models.Merchant, error)

	GetStatistic(merchantId string) (GetStatisticResponse, error)
}

func NewService() Service {
	return &service{}
}

func (s *service) GetStatistic(merchantId string) (GetStatisticResponse, error) {
	return GetStatistic(merchantId)
}

func (s *service) CreateMerchant(request models.Merchant) (models.Merchant, error) {

	if request.MerchantName == "" {
		InvalidCharacter.Message = "no merchant name"
		return models.Merchant{}, InvalidCharacter
	}
	if request.IE == "" {
		InvalidCharacter.Message = "no merchant IE"
		return models.Merchant{}, InvalidCharacter
	}
	if request.Address == "" {
		InvalidCharacter.Message = "no merchant address"
		return models.Merchant{}, InvalidCharacter
	}
	if request.BIN == "" {
		InvalidCharacter.Message = "no merchant BIN"
		return models.Merchant{}, InvalidCharacter
	}
	if request.Phone == "" {
		InvalidCharacter.Message = "no merchant phone number"
		return models.Merchant{}, InvalidCharacter
	}
	if request.Email == "" {
		InvalidCharacter.Message = "no merchant email"
		return models.Merchant{}, InvalidCharacter
	}

	merchantId := RandStringRunes(12)
	if merchantId == "" {
		return models.Merchant{}, Conflict
	}

	request.MerchantId = merchantId
	request.CreatedOn = time.Now()
	request.UpdatedOn = request.CreatedOn
	request.Status = models.MerchantStatusActive

	merchant, err := InsertMerchant(context.Background(), request)
	if err != nil {
		return merchant, InternalServerError
	}
	return merchant, err

}

func (s *service) UpdateMerchant(request models.Merchant) (models.Merchant, error) {
	var merchant models.Merchant
	if request.MerchantId == "" {
		InvalidCharacter.Message = "no merchant id"
		return models.Merchant{}, InvalidCharacter
	}
	if request.MerchantName == "" {
		InvalidCharacter.Message = "no merchant name"
		return models.Merchant{}, InvalidCharacter
	}
	if request.IE == "" {
		InvalidCharacter.Message = "no merchant IE"
		return models.Merchant{}, InvalidCharacter
	}
	if request.Address == "" {
		InvalidCharacter.Message = "no merchant address"
		return models.Merchant{}, InvalidCharacter
	}
	if request.BIN == "" {
		InvalidCharacter.Message = "no merchant BIN"
		return models.Merchant{}, InvalidCharacter
	}
	if request.Phone == "" {
		InvalidCharacter.Message = "no merchant phone number"
		return models.Merchant{}, InvalidCharacter
	}
	if request.Email == "" {
		InvalidCharacter.Message = "no merchant email"
		return models.Merchant{}, InvalidCharacter
	}

	if request.Status != models.MerchantStatusActive && request.Status != models.MerchantStatusDisabled {
		InvalidCharacter.Message = "invalid merchant status"
		return models.Merchant{}, InvalidCharacter
	}

	request.UpdatedOn = time.Now()

	err := UpdateMerchant(context.Background(), request)
	if err != nil {
		return merchant, InternalServerError
	}
	return request, nil
}

func (s *service) DeleteByIdMerchant(id string) error {

	err := DeleteMerchantById(context.Background(), id)
	if err != nil {
		return InternalServerError
	}
	return nil
}

func (s *service) DeleteBatchMerchant(ids []string) error {
	err := MDeleteMerchantByIds(context.Background(), ids)
	if err != nil {
		return InternalServerError
	}
	return nil
}

func (s *service) GetMerchantById(id string) (models.Merchant, error) {
	merchant, err := GetMerchantById(context.Background(), id)
	if err != nil {
		return merchant, InternalServerError
	}
	return merchant, nil
}

func (s *service) FilterMerchants(request FilterMerchantsRequest) ([]models.Merchant, error) {

	merchants, err := FilterMerchants(context.TODO(), request)
	if err != nil {
		return nil, Conflict
	}

	return merchants, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
