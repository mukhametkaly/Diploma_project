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
	Salt       string
	Role       string
	MerchantId string
}

func (d *UserDTO) fromDTO() (models.User, error) {
	var product models.User
	product.UserName = d.Username
	product.Password = d.Password
	product.Salt = d.Salt
	product.Role = d.Role
	product.MerchantId = d.MerchantId
	return product, nil
}

func (d *UserDTO) toDTO(product models.User) error {
	d.Username = product.UserName
	d.Password = product.Password
	d.Salt = product.Salt
	d.MerchantId = product.MerchantId
	d.Role = product.Role
	return nil
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
