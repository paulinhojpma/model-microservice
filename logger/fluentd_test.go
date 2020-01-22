package logger

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"
	"time"
)

func initFluentd() (*ILogger, error) {
	config := &OptionsConfigLogger{}
	dat, err := ioutil.ReadFile("../config-logger.json")

	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println(string(dat))
	errJSON := json.Unmarshal(dat, config)
	if errJSON != nil {
		log.Println(errJSON)
		return nil, errJSON
	}

	logger, errLo := config.ConfiguraLogger()
	if errLo != nil {
		return nil, errLo
	}
	return logger, nil
}

func TestSendLogger(t *testing.T) {
	logger, err := initFluentd()
	if err != nil {
		t.Error("Expect nothing got ", err)
		return
	}
	l := *logger
	l.Send(INFO, "testando", "11111111")
	time.Sleep(time.Second * 10)
}
