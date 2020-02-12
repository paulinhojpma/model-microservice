package cache

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	redis "github.com/go-redis/redis/v7"
)

// Redis ...
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

// GetValue ...
func (r *Redis) GetValue(key string, value interface{}) error {
	log.Println("---------Pegar objeto redis-------")
	if r.Client == nil {
		log.Println("Client redis nil")
		return errors.New("client Redis inoperante")
	}

	r.Client.Get(key).Bytes()
	conteudoRedis, erro := r.Client.Get(key).Bytes()
	if erro != nil {

		log.Println(erro)
		return erro
	}

	log.Println("Valor retornado redis -", string(conteudoRedis))
	if errj := json.Unmarshal(conteudoRedis, value); erro != nil {
		return errj
	}

	return nil
}

// SetValue ...
func (r *Redis) SetValue(key string, value interface{}, expire time.Duration) error {
	if r.Client == nil {
		log.Println("Client redis nil")
		return errors.New("client Redis inoperante")
	}
	bytString, errByte := json.Marshal(value)
	if errByte != nil {
		return errByte
	}
	erroB := r.Client.Set(key, bytString, expire).Err()
	if erroB != nil {
		return erroB
	}
	return nil
}

// DelValue ...
func (r *Redis) DelValue(key string) error {
	if r.Client == nil {
		log.Println("Client redis nil")
		return errors.New("client Redis inoperante")
	}
	error := r.Client.Del(key).Err()
	if error != nil {
		return error
	}
	return nil
}

// DelAll ...
func (r *Redis) DelAll() error {
	if r.Client == nil {
		log.Println("Client redis nil")
		return errors.New("client Redis inoperante")
	}
	err := r.Client.FlushAll().Err()
	if err != nil {
		return err
	}
	return nil
}

// AddValues ...
func (r *Redis) AddValues(key string, values interface{}) error {
	if r.Client == nil {
		log.Println("Client redis nil")
		return errors.New("client Redis inoperante")
	}

	// bytes := make([]byte, 0)
	// for i, v := range values {
	// 	byt, err := json.Marshal(v)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	log.Println(string(byt))
	// 	bytes = append(bytes, byt)
	// }
	bytes, err := json.Marshal(values)
	if err != nil {
		return err
	}
	s := string(bytes)
	if s[0:1] == "[" {
		s = s[1:]
	}
	if s[len(s)-1:] == "]" {
		s = s[:len(s)-1]
	}

	log.Println("VALOR A INSERIR - ", s)
	errRe := r.Client.LPush(key, s).Err()
	if errRe != nil {
		log.Println("ERRO AO INSERIR NA LISTA")
		return errRe
	}
	return nil
}

func (r *Redis) GetListValues(key string, values interface{}) error {
	if r.Client == nil {
		log.Println("Client redis nil")
		return errors.New("client Redis inoperante")
	}

	list, errLIs := r.Client.LRange(key, 0, -1).Result()
	if errLIs != nil {
		return errLIs
	}

	errJson := json.Unmarshal([]byte("["+strings.Join(list, ", ")+"]"), values)
	if errJson != nil {
		return errJson
	}
	// for _, v := range list {
	//
	// 	errJson := json.Unmarshal([]byte(v), values)
	// 	if errJson != nil {
	// 		return errJson
	// 	}
	// log.Println(values.(string))

	return nil
}
