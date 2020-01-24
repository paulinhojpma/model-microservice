package storage

// IStorage ...
type IStorage interface {
	connectServiceStorage(o *OptionsConfigStorage) error
	SaveFileStorage(body string, bucket string, path string) (string, error)
	GetUrlFile(bucket, path, fileName string) (string, error)
}

// OptionsConfigStorage ...
type OptionsConfigStorage struct {
	URL      string                 `json:"url"`
	Host     string                 `json:"host"`
	Password string                 `json:"password"`
	User     string                 `json:"user"`
	Driver   string                 `json:"driver"`
	Args     map[string]interface{} `json:"args"`
}

// ConfigureStorage ...
func (o *OptionsConfigStorage) ConfigureStorage() (*IStorage, error) {
	var client IStorage
	switch o.Driver {
	case "minio":
		minio := &Minio{}
		errMinio := minio.connectServiceStorage(o)
		if errMinio != nil {
			return nil, errMinio
		}
		client = minio

	}
	return &client, nil

}
