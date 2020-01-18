package cache

// ICacheClient ..
type ICacheClient interface {
	connectService(config *OptionsCacheClient) error
	// GetValues(key string, value *interface{}) error
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
