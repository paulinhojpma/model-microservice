package storage

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strings"
	"time"

	// "os"

	minio "github.com/minio/minio-go/v6"
)

// Minio ...
type Minio struct {
	Client *minio.Client
	Bucket string
}

func (m *Minio) connectServiceStorage(o *OptionsConfigStorage) error {
	minioClient, err := minio.New(o.Host, o.User, o.Password, false)
	if err != nil {
		return err
	}
	m.Client = minioClient
	m.Bucket = o.Args["bucket"].(string)
	return nil
}

// SaveFileStorage ...
func (m *Minio) SaveFileStorage(body, bucket, path string) (string, error) {
	fileB, size, fileType, errFile := base64ToByteFile(body)
	if size > 100000 {
		return "", ErrMaxFileSize
	}
	if errFile != nil {
		return "", errFile
	}
	nameFile := generateUniqueName() //+ "." + strings.Split(fileType, "/")[1]
	log.Println("nome gerado -", nameFile)
	tmpfile, err := ioutil.TempFile("", nameFile)
	if err != nil {
		log.Fatal(err)

		return "", err
	}
	log.Println("Nome do arquivo - ", tmpfile.Name())
	// defer os.Remove(tmpfile.Name()) // clean up

	if _, errTmp := tmpfile.Write(fileB); err != nil {
		log.Fatal(err)
		return "", errTmp
	}

	if path != "" {
		path = path + "/"
	}
	n, err := m.Client.FPutObject(bucket, path+nameFile+"."+strings.Split(fileType, "/")[1], tmpfile.Name(), minio.PutObjectOptions{ContentType: fileType})
	if err != nil {
		log.Fatalln(err)
	}

	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}
	log.Println("Tamanho do arquivo salvo -", n)
	return path + nameFile + "." + strings.Split(fileType, "/")[1], nil
}

// GetUrlFile ...
func (m *Minio) GetUrlFile(bucket, path, fileName string) (string, error) {
	// Set request parameters for content-disposition.
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf("inline; filename=\"%s\"", fileName))

	// Generates a presigned url which expires in a day.
	presignedURL, err := m.Client.PresignedGetObject(bucket, path, time.Second*24*60*60, reqParams)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println("Successfully generated presigned URL", presignedURL)
	return presignedURL.String(), nil
}

// CreateStorage ...
func (m *Minio) CreateStorage(bucketName string) error {
	bucketName = strings.ToLower(bucketName)
	err := m.Client.MakeBucket(bucketName, "us-east-1")
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Successfully created mybucket.")
	return nil
}

//
