package cache

import (
	"log"

	redis "github.com/go-redis/redis/v7"
)

type Redis struct {
	Client *redis.Client
}

func (r *Redis) connectService(config *OptionsCacheClient) error {
	clientRedis := redis.NewClient(&redis.Options{
		Addr:     config.Host,
		Password: config.Password,
		DB:       config.Args["DB"].(int), // use default DB
	})

	res, err := clientRedis.Ping().Result()
	// clientRedis, errRed := redis.DialURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Println("Erro ao iniciar o cache redis. Erro=", err)
		return err
		// panic("Cache Redis não foi iniciado com sucesso.")
	}
	log.Println("Redis Conectado - ", res)
	r.Client = clientRedis
	return nil
}

func (r *Redis) GetValues(key string, value *interface{}) error {
	return nil
}
