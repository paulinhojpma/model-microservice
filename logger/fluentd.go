package logger

import (
	"log"
	"time"

	"github.com/fluent/fluent-logger-golang/fluent"
)

const (
	DEBBUGER = ".debug"
	INFO     = ".info"
	WARNING  = ".warning"
	ERROR    = ".error"
	TRACE    = ".trace"
	FATAL    = ".fatal"
)

// Fluentd ...
type Fluentd struct {
	Logger  *fluent.Fluent
	Service string
}

func (f *Fluentd) connectServiceLogger(o *OptionsConfigLogger) error {
	config := fluent.Config{
		FluentPort: o.Port,
		FluentHost: o.Host,
	}
	logger, errLogger := fluent.New(config)
	if errLogger != nil {
		return errLogger
	}
	f.Logger = logger
	f.Service = o.Args["service"].(string)
	return nil
}

// Send ...
func (f *Fluentd) Send(tyype string, msg string, idOperation string) {
	go func() {
		t := time.Now()
		data := LoggerData{
			Message:     msg,
			IDOperation: idOperation,
		}
		errorPost := f.Logger.PostWithTime(f.Service+tyype, t, data)
		if errorPost != nil {
			log.Println(errorPost)
		}

	}()
}
