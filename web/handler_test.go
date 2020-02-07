package web

import (
	"fmt"
	"log"

	"sab.io/escola-service/database"
)

func connectDB() (*database.DataBase, error) {
	optionDB := &database.OptionsDB{DriverName: "postgres", IP: "tuffi.db.elephantsql.com", Porta: 5432,
		NomeDB: "rnuhlodj", User: "rnuhlodj", Senha: "adKEd_6EHdT1BV42rL_9FYJdBeCJJfmx", Debug: false, Alias: "rnuhlodj",
		TamPoolIdleConn: 1, TempoPoolIdleConn: 1, LogMinDuration: 100}
	// optionDB := &OptionsDB{DriverName: "postgres", IP: "localhost", Porta: 5432,
	// 	NomeDB: "sabio", User: "postgres", Senha: "postgres", Debug: false, Alias: "postgres",
	// 	TamPoolIdleConn: 1, TempoPoolIdleConn: 1, LogMinDuration: 100}
	db := database.NewDB(optionDB)

	if err := db.Open(); err != nil {
		log.Println("Erro ao conectar no DB. Erro =", err)
		return nil, err
	} else {
		fmt.Printf("Conectado DB OK!.\n")
	}
	return db, nil
}

func initHandler() *Handler {
	handler := &Handler{}
	DB, errDB := connectDB()
	if errDB != nil {
		log.Println(errDB)
	}
	handler.DB = DB
	return handler
}

// func TestCadastrarEscola(t *testing.T) {
// 	handler := initHandler()
// 	unidades := make([]*Unidade, 0)
// 	unidade := &Unidade{
// 		Nome: "Unidade Bancários",
// 		Endereco: &Endereco{
// 			Logradouro:  "Rua José Firmino Ferreira",
// 			Numero:      "767",
// 			Bairro:      "Agua Fria",
// 			Complemento: "Apt 303 bl A",
// 			UF:          "PB",
// 			Cidade:      "João Pessoa",
// 			Cep:         "58053222",
// 		},
// 	}
// 	unidades = append(unidades, unidade)
// 	escola := &Escola{
// 		Nome:     "Nova Escola",
// 		Cnpj:     "83164077000124",
// 		Unidades: unidades,
// 	}
// 	log.Println("Unidade inserida - ", escola.Unidades[0])
// 	log.Println("Endereco inserida - ", escola.Unidades[0].Endereco.Logradouro)
// 	errCad := escola.CadastrarEscola(handler, nil)
// 	if errCad != nil {
// 		log.Println(errCad)
// 	}
// 	if escola.IDEscola == 0 {
// 		t.Error("Expecting id greater than 0 got id -", escola.IDEscola)
// 	}
//
// }
// func TestGetEscolas(t *testing.T) {
// 	handler := initHandler()
//
// 	escolas, errEscolas := GetEscolas(handler)
// 	for _, escola := range escolas {
// 		log.Printf("%+v", escola)
// 		if escola.Unidades != nil {
// 			for _, unidade := range escola.Unidades {
// 				log.Printf("%+v", unidade)
// 			}
// 		}
// 	}
//
// 	if errEscolas != nil {
// 		t.Error("Expect nothing, got ", errEscolas)
// 	}
//
// }

// func TestGetEscola(t *testing.T) {
// 	handler := initHandler()
//
// 	escola, errEscolas := GetEscola(handler, 16)
// 	log.Printf("%+v", escola)
// 	log.Printf("%+v", escola.Unidades[0])
//
// 	if errEscolas != nil {
// 		t.Error("Expect nothing, got ", errEscolas)
// 	}
// }

// func initHandler() *web.Handler {
// 	handler := &web.Handler{}
// 	messaging := configMessage()
// 	handler.Messaging = messaging
// 	return handler
// }
// func configMessage() *messaging.IMessageClient {
// 	config := initializeConfigMessage()
// 	IMessage, _ := config.ConfiguraFilaMensagens()
//
// 	return IMessage
// }
// func initializeConfigMessage() *messaging.OptionsMessageCLient {
// 	iMessage := &messaging.OptionsMessageCLient{}
//
// 	dat, err := ioutil.ReadFile("../config-msgs.json")
//
// 	if err != nil {
// 		log.Println(err)
// 		return nil
// 	}
// 	log.Println(string(dat))
// 	errJSON := json.Unmarshal(dat, iMessage)
// 	if errJSON != nil {
// 		log.Println(errJSON)
// 		return nil
// 	}
// 	log.Println(iMessage.Args)
// 	//iMessage.ConfiguraFilaMensagens
// 	return iMessage
//
// }

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

// func TestSendMessage(t *testing.T) {
//
// 	handler := initHandler()
// 	param := &messaging.MessageParam{}
// 	param.Method = "GET"
// 	param.Params = map[string]int{
// 		"escola":  1,
// 		"unidade": 2,
// 	}
// 	param.Resource = "unidade"
// 	param.Type = "request"
// 	param.Status = 0
// 	param.Body = []byte("VICENTE")
// 	messag := *handler.Messaging
// 	error := messag.PublishMessage("escola", param)
// 	if error != nil {
// 		t.Error("Expected message publishe, got ", error.Error())
// 	}
// }

// func TestReceiveMessage(t *testing.T) {
// 	handler := initHandler()
// 	messag := *handler.Messaging
// 	msgChan, errMsg := messag.ReceiveMessage("escola")
// 	if errMsg != nil {
// 		log.Println(errMsg)
// 	}
// 	msg := &messaging.MessageParam{}
// 	for m := range msgChan {
// 		msg = &m
// 		break
// 	}
// 	fmt.Printf("%+v\n", msg)
// 	if msg == nil {
// 		t.Error("Expected message exist, got ", nil)
// 	}
// }

// func TestSendAndReceiveMessage(t *testing.T) {
// 	handler := initHandler()
// 	param := &messaging.MessageParam{}
// 	param.Method = "GET"
// 	param.Params = map[string]int{
// 		"escola":  1,
// 		"unidade": 2,
// 	}
// 	param.Resource = "unidade"
// 	param.Type = "request"
// 	param.Status = 0
// 	param.Body = []byte("VICENTE")
// 	messag := *handler.Messaging
// 	msg, errMsg := messag.ReceiveMessage("escola")
// 	if errMsg != nil {
// 		log.Println("Erro de recebimento - ", errMsg)
// 	}
// 	go func() {
// 		for d := range msg {
// 			log.Println("Nome da fila a ser enviada - ", d.Args["replyTo"].(string))
// 			log.Println("Mensagem a ser enviada - ", string(d.Body))
// 			err := messag.RespondMessage(d.Args["replyTo"].(string), &d)
// 			if err != nil {
// 				log.Println(err)
// 			}
// 			log.Println("Respondendo mensagem recebida")
// 		}
// 	}()
//
// 	msgRecebida, error := messag.PublishAndReceiveMessage("escola", param)
// 	log.Println("Mensagem Recebida - ", string(msgRecebida.Body))
// 	if error != nil {
//
// 		t.Error("Expected message publishe, got ", error)
// 	}
// }

// func TestGetDisciplinas(t *testing.T) {
// 	handler := initHandler()
//
// 	disciplinas, errDisciplinas := GetDisciplinas(handler, 16)
// 	if disciplinas != nil {
// 		for _, disciplina := range disciplinas {
// 			fmt.Printf("%+v\n", disciplina)
// 			for _, serie := range disciplina.Ementas {
// 				fmt.Printf("%+v\n", serie)
// 				fmt.Printf("%+v\n", serie.Serie)
// 			}
// 		}
// 	}
// 	if errDisciplinas != nil {
// 		t.Error("expecting nothing got, ", errDisciplinas)
// 	}
//
// }

// func TestGetDisciplinaByID(t *testing.T) {
// 	handler := initHandler()
//
// 	disciplina, errDisciplinas := GetDisciplinaByID(handler, 16, 1)
//
// 	if disciplina != nil {
// 		fmt.Printf("%+v\n", disciplina)
// 		for _, serie := range disciplina.Ementas {
// 			fmt.Printf("%+v\n", serie)
// 			fmt.Printf("%+v\n", serie.Serie)
// 		}
// 	}
//
// 	if errDisciplinas != nil {
// 		t.Error("expecting nothing got, ", errDisciplinas)
// 	}
//
// }

// func TestCadastrarDisciplina(t *testing.T) {
// 	handler := initHandler()
// 	sDisciplina := `{
//
//    "nome":"Português",
//    "descricao":"Disciplina de português",
//    "ementas":[
//       {
//
//          "ementa":"Morfologia",
//          "cargaHoraria":4,
//          "serie":{
//             "idSerie":1,
//             "nome":"1",
//             "tipo":"ENSINO MEDIO"
//          }
//       },
//       {
//
//          "ementa":"Sintaxe",
//          "cargaHoraria":3,
//          "serie":{
//             "idSerie":2,
//             "nome":"2",
//             "tipo":"ENSINO MEDIO"
//          }
//       },
//       {
//
//          "ementa":"Semantica",
//          "cargaHoraria":6,
//          "serie":{
//             "idSerie":3,
//             "nome":"3",
//             "tipo":"ENSINO MEDIO"
//          }
//       }
//    ]
// }`
//
// 	disciplina := &Disciplina{}
// 	errJSON := json.Unmarshal([]byte(sDisciplina), disciplina)
// 	if errJSON != nil {
// 		log.Println(errJSON)
// 	}
// 	errDisciplina := disciplina.CadastrarDisciplina(handler, 16, nil)
// 	log.Println("ID disciplina ", disciplina.IDDisciplina)
// 	if errDisciplina != nil {
// 		t.Error("Expecting nothing got ", errDisciplina)
// 	}
// }
