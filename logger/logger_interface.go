package logger

type ILogger interface {
	connectServiceLogger(o *OptionsConfigLogger) error
	Send(tyype string, msg string, idOperation string)
}

type OptionsConfigLogger struct {
	URL    string                 `json:"url"`
	Host   string                 `json:"host"`
	Senha  string                 `json:"senha"`
	Port   int                    `json:"port"`
	Driver string                 `json:"driver"`
	Args   map[string]interface{} `json:"args"`
}

type LoggerData struct {
	Message     string `json:"message"`
	IDOperation string `json:"idOperation"`
}

func (o *OptionsConfigLogger) ConfiguraLogger() (*ILogger, error) {
	var client ILogger
	switch o.Driver {
	case "fluentd":
		fluentd := &Fluentd{}
		errFluentd := fluentd.connectServiceLogger(o)
		if errFluentd != nil {
			return nil, errFluentd
		}
		client = fluentd
	}
	return &client, nil
}
