package config

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Configs struct {
	LogrusLevel uint8          `json:"logrus_level"`
	Postgres    PostgresConfig `json:"postgres"`
}

type PostgresConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	DBName   string `json:"dbname"`
}

var AllConfigs *Configs

func GetConfigs() error {
	var filePath string
	if os.Getenv("config") == "" {
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		filePath = pwd + "/src/config/config.json"
	} else {
		filePath = os.Getenv("config")
	}
	file, err := os.Open(filePath)

	if err != nil {
		return err
	}
	decoder := json.NewDecoder(file)
	var configs Configs
	err = decoder.Decode(&configs)

	if err != nil {
		return err
	}
	AllConfigs = &configs
	return nil
}

func Healthchecks(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}
