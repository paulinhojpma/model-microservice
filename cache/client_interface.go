package cache

import "time"

// ICacheClient ..
type ICacheClient interface {
	connectService(config *OptionsCacheClient) error
	GetValue(key string, value interface{}) error
	SetValue(key string, value interface{}, expire time.Duration) error
	DelValue(key string) error
	DelAll() error
	AddValues(key string, values interface{}) error
	GetListValues(key string, values interface{}) error
}

// OptionsCacheClient ..
type OptionsCacheClient struct {
	URL      string                 `json:"url"`
	Driver   string                 `json:"driver"`
	Host     string                 `json:"host"`
	Password string                 `json:"password"`
	Args     map[string]interface{} `json:"args"`
}

// ConfiguraCache
func (o *OptionsCacheClient) ConfiguraCache() (*ICacheClient, error) {
	var client ICacheClient
	switch o.Driver {
	case "redis":
		redis := &Redis{}
		errRedis := redis.connectService(o)
		if errRedis != nil {
			return nil, errRedis
		}
		client = redis

	}
	return &client, nil
}
