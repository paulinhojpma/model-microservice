package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	// env "github.com/caarlos0/env"
)

var config *Configuracoes

//Configuracoes ...
type Configuracoes struct {
	DBBizNome      string `json:"DBBizNome" env:"DB_BIZ_NOME"`
	DBBizHost      string `json:"DBBizHost" env:"DB_BIZ_HOST"`
	DBBizPorta     int    `json:"DBBizPorta" env:"DB_BIZ_PORTA"`
	DBBizUser      string `json:"DBBizUser" env:"DB_BIZ_USER"`
	DBBizSenha     string `json:"DBBizSenha" env:"DB_BIZ_SENHA"`
	DatabaseDriver string `json:"databaseDriver" env:"DATABASE_DRIVER"`
	RedisHost      string `json:"redisHost" env:"REDIS_HOST"`
	RedisSenha     string `json:"redisSenha" env:"REDIS_SENHA"`
	RedisURL       string `json:"redisURL" env:"REDIS_URL"`
	CacheDriver    string `json:"cacheDriver" env:"CACHE_DRIVER"`
	NewRelicToken  string `json:"newRelicToken" env:"NEWRELIC_TOKEN"`
	Port           int    `json:"port" env:"PORT"`
	AllowedParam   string `json:"allowedParam" env:"ALLOWED_PARAM"`

	AWSRegion string `json:"AWSRegion" env:"AWS_REGION"`

	EmailSender     string `json:"emailSender" env:"EMAIL_SENDER"`
	EmailPassword   string `json:"emailSenha" env:"EMAIL_PASSWORD"`
	LoggerHost      string `json:"logger-host" env:"LOGGER_HOST"`           // "logger-host": "0.0.0.0",
	LoggerPort      int    `json:"logger-port" env:"LOGGER_PORT"`           // "logger-port": 9880,
	LoggerDriver    string `json:"logger-driver" env:"LOGGER_DRIVER"`       // "logger-driver": "fluentd",
	LoggerService   string `json:"logger-service" env:"LOGGER_SERVICE"`     // "logger-service": "escola",
	Messaging       string `json:"messaging" env:"MESSAGING"`               // "messaging": "../../config-msgs.json",
	StorageUser     string `json:"storage-user" env:"STORAGE_USER"`         // "storage-user": "VICENTE",
	StoragePassword string `json:"storage-password" env:"STORAGE_PASSWORD"` // "storage-password": "VICENTERONALDO",
	StorageHost     string `json:"storage-host" env:"STORAGE_HOST"`         // "storage-host": "localhost:9000",
	StorageDriver   string `json:"storage-driver" env:"STORAGE_DRIVER"`     // "storage-driver": "minio",

}

//NewConfig ...
func NewConfig(file string) *Configuracoes {
	var erro error

	conf := &Configuracoes{}

	if file != "" {
		fmt.Println(file)

		bufConf, err := ioutil.ReadFile(file)
		if err == nil {
			erro = json.Unmarshal(bufConf, conf)
			if erro != nil {
				log.Println(erro)
			}
		}
	}

	// //variaveis de ambiente sobrescrevem informacoes do json
	// if erro = env.Parse(conf); erro != nil {
	// 	log.Println(erro)
	// }

	return conf
}

//Config ...
func Config() *Configuracoes {
	return config
}
