package auth

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mukhametkaly/Diploma_project/auth-api/src/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
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
	Auth(token string) (bool, error)
}

type aclClaim struct {
	models.User
	jwt.RegisteredClaims
}

func (s service) Registration(user models.User) error {

	if user.Username == "" {
		return newErrorString(http.StatusBadRequest, "no user name")
	}
	if user.Password == "" || len(user.Password) < 8 {
		return newErrorString(http.StatusBadRequest, "no password")
	}
	if user.MerchantId == "" {
		return newErrorString(http.StatusBadRequest, "no merchant id")
	}

	_, err := GetUser(context.Background(), user.Username)
	if err != pg.ErrNoRows {
		return newErrorString(http.StatusConflict, "user with this username exists")
	}

	count, err := GetUserByMerchant(context.TODO(), user.MerchantId)
	if count != 0 {
		user.Role = Cashier
	} else {
		user.Role = Merchant
	}

	user.Salt = RandStringRunes(8)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password+"@"+user.Salt), 14)
	user.Password = string(passwordHash)

	err = RegistrateUser(context.Background(), user)
	if err != nil {
		return newError(http.StatusInternalServerError, err)
	}

	return nil
}

func (s service) LogIn(username, password string) (string, error) {

	if username == "" {
		return "", newErrorString(http.StatusBadRequest, "no username")
	}

	if password == "" {
		return "", newErrorString(http.StatusBadRequest, "no password")
	}

	user, err := GetUser(context.Background(), username)
	if err != nil {
		return "", newError(http.StatusNotFound, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password+"@"+user.Salt))
	if err != nil {
		return "", newErrorString(http.StatusBadRequest, "wrong password")
	}

	if user.Role == "" {
		return "", newErrorString(http.StatusInternalServerError, "no role user"+user.Username)
	}

	if user.MerchantId == "" && user.Role != Admin {
		return "", newErrorString(http.StatusInternalServerError, "no merchantId"+user.Username)
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
		return "", newErrorString(http.StatusInternalServerError, "undefined role "+user.Role+" user "+user.Username)
	}

	token, err := newToken(user)
	if err != nil {
		return "", newErrorString(http.StatusInternalServerError, "something went wrong")
	}

	return token, nil
}

func (s service) Auth(token string) (bool, error) {

	user := getAcl(token)
	if user.Username == "" {
		return false, newErrorString(http.StatusUnauthorized, "invalid token")
	}
	return true, nil

}
