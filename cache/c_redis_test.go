package cache

import (
	"testing"
)

func TestCacheRedisCon(t *testing.T) {
	arg := make(map[string]interface{})
	arg["DB"] = 1
	confi := OptionsCacheClient{
		Host:     "127.0.0.1:6379",
		Password: "",
		Driver:   "redis",
		Args:     arg,
	}
	_, err := confi.ConfiguraCache()

	if err != nil {
		t.Error("Expected nothing, got ", nil)
	}

}
