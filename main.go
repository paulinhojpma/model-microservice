package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"sab.io/escola-service/cache"
	"sab.io/escola-service/config"
	"sab.io/escola-service/database"
	"sab.io/escola-service/logger"
	"sab.io/escola-service/messaging"
	"sab.io/escola-service/storage"
	"sab.io/escola-service/web"
)

const (
	BUCKET        = "test"
	MESSAGECONFIG = "config-msgs.json"
)

func main() {
	// time.Sleep(time.Second * 15)
	configNew := config.Config()
	stage := os.Getenv("STAGE")
	log.Println("ENVIROMENT - ", stage)
	if stage == "" {
		configNew = config.NewConfig("conf.json")
	} else {
		configNew = config.NewConfig("")
	}
	handler := &web.Handler{}
	// log.Println(handler)
	// iniciando o cache do redis
	configCache := cache.OptionsCacheClient{
		Host:     configNew.RedisHost,
		Password: configNew.RedisSenha,
		Driver:   configNew.CacheDriver,
		Args: map[string]interface{}{
			"DB": 1,
		},
	}
	redis, errRedis := configCache.ConfiguraCache()
	if errRedis != nil {
		log.Println(errRedis)
		time.Sleep(time.Second * 5)
		for {
			redis, errRedis = configCache.ConfiguraCache()
			if errRedis == nil {
				break
			}
		}
		// os.Exit(1)
	}
	log.Println("Redis Conectado.")
	log.Println("Vicente Ronaldo Conectado.")
	handler.Cache = redis

	//inicializando o storage
	optionStorage := storage.OptionsConfigStorage{
		Host:     configNew.StorageHost,
		User:     configNew.StorageUser,
		Password: configNew.StoragePassword,
		Driver:   configNew.StorageDriver,
		Args: map[string]interface{}{
			"bucket": BUCKET,
		},
	}

	clientStorage, errStorage := optionStorage.ConfigureStorage()
	if errStorage != nil {
		log.Println("Erro ao conectar ao serviço de storage - ", errStorage)
		for {
			time.Sleep(time.Second * 5)
			clientStorage, errStorage = optionStorage.ConfigureStorage()
			if errStorage == nil {
				break
			}
			log.Println("Erro ao conectar ao serviço de storage - ", errStorage)
		}
	}
	log.Println("Serviço de storage conectado")
	handler.Storage = clientStorage
	// inicializand o o client Logger
	configLogger := logger.OptionsConfigLogger{
		Host:   configNew.LoggerHost,
		Port:   configNew.LoggerPort,
		Driver: configNew.LoggerDriver,
		Args: map[string]interface{}{
			"service": configNew.LoggerService,
		},
	}
	log.Println("HOST RETORNADO - ", configNew.LoggerHost)
	clientLogger, errLogger := configLogger.ConfiguraLogger()
	if errLogger != nil {
		log.Println("ERRO ao conectar ao serviço de logger -", errLogger)
		for {
			time.Sleep(time.Second * 5)
			clientLogger, errLogger = configLogger.ConfiguraLogger()
			if errLogger == nil {
				break
			}
			log.Println("ERRO ao conectar ao serviço de logger -", errLogger)
		}
	}
	log.Println("Logger conectado")

	handler.Logger = clientLogger

	// loggerr := *clientLogger

	// clientLogger.Send(logger.INFO, "testando", "1111111")
	(*handler.Logger).Send(logger.INFO, "testando", "1111111")

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
	configMessage := &messaging.OptionsMessageCLient{
		URL:    configNew.MessagingURL,
		Driver: configNew.MessagingDriver,
	}

	dat, errRe := ioutil.ReadFile(MESSAGECONFIG)
	if errRe != nil {
		log.Println(errRe)

	}
	errJSON := json.Unmarshal(dat, configMessage)
	if errJSON != nil {
		log.Println(errJSON)
	}

	log.Println(configMessage.URL)
	log.Println(configMessage.Driver)
	log.Println(configMessage.Args)
	imessa, errM := configMessage.ConfiguraFilaMensagens()

	if errM != nil {
		log.Println("Erro ao conectar na fila de mensagens - ", errM)
		log.Println("Antes do for")
		for {

			time.Sleep(time.Second * 5)
			imessa, errM = configMessage.ConfiguraFilaMensagens()
			log.Println(errM)
			if errM == nil {
				break
			}
		}

	}
	log.Println("Conectado ao serviço de menssagens")

	// messa := *imessa
	handler.Message = imessa
	msgChan, errMsg := (*imessa).ReceiveMessage("escola")
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
			path := web.PathMethod(msg.Method, msg.Resource)
			newMsg := manager.CallService(path, msg)
			log.Println(newMsg)
			errMessa := manager.ManagerMessage(*imessa, newMsg)
			if errMessa != nil {
				// TODO: colocar o serviço de loggerr
				log.Println("erro a retornar a mensagem - ", errMsg)
			}
		}
	}()

	<-forever
}
