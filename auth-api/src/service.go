package src

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mukhametkaly/Diploma/auth-api/models"
	"golang.org/x/crypto/bcrypt"
)

const (
	Admin    = "admin"
	Merchant = "merchant"
	Cashier  = "cashier"
)

type service struct {
}

func NewService() Service {
	return &service{}
}

// Service is the interface that provides methods.
type Service interface {
	LogIn(username, password string) (string, error)
	Registration(user models.User) error
	Auth(token string) (models.Rights, error)
}

type aclClaim struct {
	models.User
	jwt.RegisteredClaims
}

func (s service) Registration(user models.User) error {

	if user.Username == "" {
		return nil
	}
	if user.Password == "" || len(user.Password) < 8 {
		return nil
	}
	if user.MerchantId == "" {
		return nil
	}
	if user.Role == Admin || user.Role == Merchant || user.Role == Cashier {
		return nil
	}

	user, err := GetUser(context.Background(), user.Username)
	if err == nil && user.Password != "" {
		panic(err)
	}

	user.Salt = RandStringRunes(8)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password+"@"+user.Salt), 14)
	user.Password = string(passwordHash)

	err = RegistrateUser(context.Background(), user)
	if err != nil {
		panic(err)
	}

	return nil
}

func (s service) LogIn(username, password string) (string, error) {

	user, err := GetUser(context.Background(), username)
	if err != nil {
		return "", err
	}

	reqHash, err := bcrypt.GenerateFromPassword([]byte(password+"@"+user.Salt), 14)

	err = bcrypt.CompareHashAndPassword(reqHash, []byte(user.Password+"@"+user.Salt))
	if err != nil {
		return "", err
	}

	if user.Role == "" {
		return "", err
	}

	if user.MerchantId == "" {
		return "", err
	}

	acl := make(map[string]models.Rights)
	switch user.Role {
	case Admin:
		acl[Admin] = AdminRights
	case Merchant:
		acl[user.MerchantId] = MerchantRights
	case Cashier:
		acl[user.MerchantId] = CashierRights
	default:
		return dfs, ""
	}

	return "", nil
}

func (s service) Auth(token string) (models.Rights, error) {

	tk, err := jwt.ParseWithClaims(token, &aclClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("AllYourBase"), nil
	})
	if err != nil {
		return models.Rights{}, errors.New("internal error")
	}

	if claim, ok := tk.Claims.(*aclClaim); ok && tk.Valid {
		return claim.Rights, nil
	} else {
		return models.Rights{}, errors.New("unauthorized")
	}

}
