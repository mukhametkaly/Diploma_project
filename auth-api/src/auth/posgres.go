package auth

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/mukhametkaly/Diploma_project/auth-api/src/config"
	"github.com/mukhametkaly/Diploma_project/auth-api/src/models"
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
	tableName  struct{} `pg:"users"`
	Username   string   `pg:",pk,unique"`
	Password   string
	FullName   string
	IIN        string
	Mail       string
	Mobile     string
	Salt       string
	UserRole   string
	MerchantId string `pg:"merchant_id"`
}

func (d *UserDTO) fromDTO() models.User {
	var user models.User
	user.IIN = d.IIN
	user.Mail = d.Mail
	user.FullName = d.FullName
	user.Username = d.Username
	user.Password = d.Password
	user.Mobile = d.Mobile
	user.Salt = d.Salt
	user.Role = d.UserRole
	user.MerchantId = d.MerchantId
	return user
}

func (d *UserDTO) toDTO(user models.User) {
	d.Username = user.Username
	d.Password = user.Password
	d.FullName = user.FullName
	d.IIN = user.IIN
	d.Mail = user.Mail
	d.Mobile = user.Mobile
	d.Salt = user.Salt
	d.MerchantId = user.MerchantId
	d.UserRole = user.Role
}

func GetUser(ctx context.Context, username string) (user models.User, err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in GetUser", err.Error())
		return
	}

	userDto := UserDTO{}

	//todo where username and password
	q := conn.ModelContext(ctx, &userDto).Where("username = ?", username)
	err = q.Select()
	if err != nil {
		Loger.Debugln("error select in get list users", err.Error())
		return
	}

	return userDto.fromDTO(), nil
}

func RegistrateUser(ctx context.Context, user models.User) (err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in RegistrateUser", err.Error())
		return
	}

	userDto := UserDTO{}
	userDto.toDTO(user)

	_, err = conn.ModelContext(ctx, &userDto).Returning("*", &userDto).Insert()
	if err != nil {
		Loger.Debugln("error insert", err.Error())
		return
	}

	return
}

func GetUserByMerchant(ctx context.Context, merchantId string) (count int, err error) {
	conn, err := GetPGSession()
	if err != nil {
		Loger.Debugln("error getSession in GetUserByMerchant", err.Error())
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
