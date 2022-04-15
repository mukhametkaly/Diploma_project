package src

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/mukhametkaly/Diploma/auth-api/config"
	"github.com/mukhametkaly/Diploma/auth-api/models"
	"time"

	"github.com/sirupsen/logrus"
)

var db *pg.DB

func PGConnectStart() (*pg.DB, error) {
	conn := pg.Connect(&pg.Options{
		Addr:               fmt.Sprintf("%s:%s", config.AllConfigs.Postgres.Host, config.AllConfigs.Postgres.Port),
		User:               config.AllConfigs.Postgres.User,
		Password:           config.AllConfigs.Postgres.Password,
		Database:           config.AllConfigs.Postgres.DBName,
		IdleTimeout:        59 * time.Second,
		IdleCheckFrequency: 30 * time.Second,
	})

	err := conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func GetPGSession() (*pg.DB, error) {

	if db == nil {
		client, err := PGConnectStart()
		if err != nil {
			return nil, err
		} else {
			db = client
			return db, nil
		}
	} else {
		return db, nil
	}
}

// Repository is postgres repository
type Repository struct {
	db     *pg.DB
	logger logrus.Logger
}

type UserDTO struct {
	tableName  struct{} `pg:"user"`
	Username   string   `pg:",pk,unique"`
	Password   string
	FullName   string
	IIN        string
	Mail       string
	Salt       string
	Role       string
	MerchantId string
}

func (d *UserDTO) fromDTO() models.User {
	var user models.User
	user.IIN = d.IIN
	user.Mail = d.Mail
	user.FullName = d.FullName
	user.Username = d.Username
	user.Password = d.Password
	user.Salt = d.Salt
	user.Role = d.Role
	user.MerchantId = d.MerchantId
	return user
}

func (d *UserDTO) toDTO(product models.User) {
	d.Username = product.Username
	d.Password = product.Password
	d.FullName = product.FullName
	d.IIN = product.IIN
	d.Mail = product.Mail
	d.Salt = product.Salt
	d.MerchantId = product.MerchantId
	d.Role = product.Role
}

func GetUser(ctx context.Context, username string) (user models.User, err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in GetProductById", err.Error())
		return
	}

	//todo where username and password
	q := conn.ModelContext(ctx, &user).Where("username = ?", username)
	err = q.Select()
	if err != nil {
		Loger.Debugln("error select in get list orders", err.Error())
		return
	}

	return
}

func RegistrateUser(ctx context.Context, user models.User) (err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in GetProductById", err.Error())
		return
	}

	_, err = conn.ModelContext(ctx, &user).Insert(&user)
	if err != nil {
		Loger.Debugln("error select in get list orders", err.Error())
		return
	}

	return
}

func GetUserByMerchant(ctx context.Context, merchantId string) (count int, err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in GetProductById", err.Error())
		return
	}

	userDto := UserDTO{}

	//todo where username and password
	count, err = conn.ModelContext(ctx, &userDto).Where("merchant_id = ?", merchantId).Count()
	if err != nil {
		Loger.Debugln("error count in get list users", err.Error())
		return
	}

	return

}
