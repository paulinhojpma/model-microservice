package main

import (
	"github.com/sab.io/escola-service/cache"
)

func main() {
	config := cache.OptionsCacheClient{
		Host:     "127.0.0.1:6379",
		Password: "",
		Driver:   "redis",
	}
	_, err := config.ConfiguraCache()
	if err != nil {
		return
	}

}
