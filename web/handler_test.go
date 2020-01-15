package web_test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"

	messaging "github.com/sab.io/escola-service/messaging"
	"github.com/sab.io/escola-service/web"
)

func initHandler() *web.Handler {
	handler := &web.Handler{}
	messaging := configMessage()
	handler.Messaging = messaging
	return handler
}
func configMessage() *messaging.IMessageClient {
	config := initializeConfigMessage()
	IMessage, _ := config.ConfiguraFilaMensagens()

	return IMessage
}
func initializeConfigMessage() *messaging.OptionsMessageCLient {
	iMessage := &messaging.OptionsMessageCLient{}

	dat, err := ioutil.ReadFile("../config-msgs.json")

	if err != nil {
		log.Println(err)
		return nil
	}
	log.Println(string(dat))
	errJSON := json.Unmarshal(dat, iMessage)
	if errJSON != nil {
		log.Println(errJSON)
		return nil
	}
	log.Println(iMessage.Args)
	//iMessage.ConfiguraFilaMensagens
	return iMessage

}

// func TestConfigMessage(t *testing.T) {
// 	config := initializeConfigMessage()
// 	imessa, error := config.ConfiguraFilaMensagens()
// 	if error != nil {
//
// 		log.Println(error)
//
// 	}
//
// 	messa := *imessa
//
// 	rab := messa.(*messaging.Rabbit)
// 	log.Println(rab.Exchange.Binding["turma"][0].Name)
// 	if _, ok := rab.Exchange.Binding["escola"]; !ok {
// 		t.Error("Expected escola biding, got ", "nothing")
// 	}
//
// }

func TestSendMessage(t *testing.T) {
	handler := initHandler()
	param := &messaging.MessageParam{}
	param.Method = "GET"
	param.Params = map[string]int{
		"escola":  1,
		"unidade": 2,
	}
	param.Resource = "unidade"
	param.Type = "request"
	param.Status = 0
	messag := *handler.Messaging
	error := messag.PublishMessage("escola", param)
	if error != nil {
		t.Error("Expected message publishe, got ", error.Error())
	}
}
