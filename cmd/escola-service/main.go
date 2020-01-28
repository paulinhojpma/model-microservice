package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"sab.io/escola-service/config"
	"sab.io/escola-service/database"
	"sab.io/escola-service/messaging"
	"sab.io/escola-service/web"
)

func main() {
	configNew := config.Config()
	stage := os.Getenv("stage")
	if stage == "" {
		configNew = config.NewConfig("conf.json")
	} else {
		configNew = config.NewConfig("")
	}
	handler := &web.Handler{}
	log.Println(handler)
	// // iniciando o cache do redis
	// configCache := cache.OptionsCacheClient{
	// 	Host:     configNew.RedisHost,
	// 	Password: configNew.RedisSenha,
	// 	Driver:   configNew.CacheDriver,
	// 	Args: map[string]interface{}{
	// 		"DB": 1,
	// 	},
	// }
	// redis, errRedis := configCache.ConfiguraCache()
	// if errRedis != nil {
	// 	log.Println(errRedis)
	// 	os.Exit(1)
	// }
	// log.Println("Redis Conectado.")
	// log.Println("Vicente Ronaldo Conectado.")
	// handler.Cache = *redis
	//
	// // inicializando o client Logger
	// configLogger := logger.OptionsConfigLogger{
	// 	URL:    configNew.LoggerHost,
	// 	Port:   configNew.LoggerPort,
	// 	Driver: configNew.LoggerDriver,
	// 	Args: map[string]interface{}{
	// 		"service": configNew.LoggerService,
	// 	},
	// }
	// clientLogger, errLogger := configLogger.ConfiguraLogger()
	// if errLogger != nil {
	// 	log.Println("ERRO ao conectar ao serviço de logger -", errLogger)
	// }
	// log.Println("Logger conectado")
	//
	// handler.Logger = *clientLogger
	// handler.Logger.Send(logger.INFO, "testando", "1111111")

	// inicializando serviço do DataBase
	optionDB := &database.OptionsDB{DriverName: configNew.DatabaseDriver, IP: configNew.DBBizHost, Porta: configNew.DBBizPorta,
		NomeDB: configNew.DBBizNome, User: configNew.DBBizUser, Senha: configNew.DBBizSenha, Debug: false, Alias: configNew.DBBizNome,
		TamPoolIdleConn: 1, TempoPoolIdleConn: 1, LogMinDuration: 100}

	db := database.NewDB(optionDB)

	if err := db.Open(); err != nil {
		log.Println("Erro ao conectar no DB. Erro=", err)
	}
	log.Println("Conectado DB OK!")
	handler.DB = db

	//inicializando serviço de messaging
	configMessage := &messaging.OptionsMessageCLient{}

	dat, err := ioutil.ReadFile(configNew.Messaging)
	if err != nil {
		log.Println(err)

	}
	errJSON := json.Unmarshal(dat, configMessage)
	if errJSON != nil {
		log.Println(errJSON)
	}

	imessa, error := configMessage.ConfiguraFilaMensagens()
	if error != nil {
		log.Println(error)
	}
	log.Println("Conectado ao serviço de menssagens")
	messa := *imessa
	handler.Message = messa
	msgChan, errMsg := handler.Message.ReceiveMessage("escola")
	if errMsg != nil {
		log.Println(errMsg)
	}
	forever := make(chan bool)
	manager := web.NewRoutes(handler)
	go func() {
		for m := range msgChan {
			msg := &m
			log.Println("recebendo mensagem")
			log.Println(string(msg.Body))
			log.Println("Método - ", msg.Method)
			newMsg := manager.CallService("test", msg)
			log.Println(newMsg)

		}
	}()

	<-forever
}
