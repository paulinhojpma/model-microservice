package storage

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"
)

func initStorage() (*IStorage, error) {
	config := &OptionsConfigStorage{}
	dat, err := ioutil.ReadFile("../config-storage.json")
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
	storage, errStorage := config.ConfigureStorage()
	if errStorage != nil {
		return nil, errStorage
	}
	return storage, nil

}

// func TestSaveFile(t *testing.T) {
// 	s, errStorage := initStorage()
// 	if errStorage != nil {
// 		t.Error("Expect nothing, got ", errStorage)
// 	}
// 	storage := *s
//
// 	dat, err := ioutil.ReadFile("../000-18960-254.PDF")
// 	if err != nil {
// 		t.Error("Expect nothing, got ", err)
// 	}
// 	dat64 := b64.StdEncoding.EncodeToString(dat)
// 	log.Println("Base64 - ", string(dat64[0:5]))
// 	_, errSave := storage.SaveFileStorage("data:application/pdf;base64,"+dat64, "teste", "klebernilton")
// 	if errSave != nil {
// 		t.Error("Expect nothing, got ", errSave)
// 	}
// }

func TestGetURLFile(t *testing.T) {
	s, errStorage := initStorage()
	if errStorage != nil {
		t.Error("Expect nothing, got ", errStorage)
	}
	storage := *s
	urls, errURL := storage.GetUrlFile("teste", "vicente.pdf", "ronaldo.pdf")
	log.Println("URL retornada - ", urls)
	if errURL != nil {
		t.Error("Expect nothing got ", errURL)
	}
}

func TestCreateStorage(t *testing.T) {
	s, errStorage := initStorage()

	if errStorage != nil {
		t.Error("Expect nothing, got ", errStorage)
	}
	storage := *s
	errBuck := storage.CreateStorage("novobucket")
	if errBuck != nil {
		t.Error("Expect nothing got ", errBuck)
	}
}
