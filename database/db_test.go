package database

import (
	"fmt"
	"log"
	"testing"

	"github.com/fluent/fluent-logger-golang/fluent"
)

func connectDB() (*DataBase, error) {
	optionDB := &OptionsDB{DriverName: "postgres", IP: "tuffi.db.elephantsql.com", Porta: 5432,
		NomeDB: "rnuhlodj", User: "rnuhlodj", Senha: "5SkFel9p6deAmbbtHyo9I-_vOwX1v-ww", Debug: false, Alias: "rnuhlodj",
		TamPoolIdleConn: 1, TempoPoolIdleConn: 1, LogMinDuration: 100}
	// optionDB := &OptionsDB{DriverName: "postgres", IP: "localhost", Porta: 5432,
	// 	NomeDB: "sabio", User: "postgres", Senha: "postgres", Debug: false, Alias: "postgres",
	// 	TamPoolIdleConn: 1, TempoPoolIdleConn: 1, LogMinDuration: 100}
	db := NewDB(optionDB)

	if err := db.Open(); err != nil {
		log.Println("Erro ao conectar no DB. Erro=", err)
		return nil, err
	} else {
		fmt.Printf("Conectado DB OK!.\n")

	}
	return db, nil
}

// func TestConnectDB(t *testing.T) {
// 	_, erro := connectDB()
// 	if erro != nil {
// 		t.Error("Expected nothing, got ", erro)
// 	}
// }

type Escola struct {
	IDEscola int64      `json:"id_escola"`
	Nome     string     `json:"nome"`
	Unidades []*Unidade `json:"unidades"`
}

type Unidade struct {
	IDUnidade int64  `json:"id_unidade"`
	Nome      string `json:"nome"`
}

// func TestQuery(t *testing.T) {
// 	db, errConn := connectDB()
// 	if errConn != nil {
// 		log.Println("ERRO CONEXÂO - ", errConn)
// 	}
// 	sql := `select row_to_json(t)
// 	        from (
// 	        	select e.id_escola, e.nome ,
// 	        		(select array_to_json(array_agg(row_to_json(d)))
// 	        			from(
// 	        				select u.id_unidade, u.nome
// 	        				from escola.unidade u
// 	        				where u.id_escola = %(id_escola)d
// 	        			) d
// 	        		) as unidades
// 	        	from escola.escola e
// 	        where e.id_escola = %(id_escola)d
// 	        ) t;`
// 	// sql := `select e.nome from escola.escola e where e.id_escola = 1; `
// 	argMap := map[string]interface{}{
// 		"id_escola": 1,
// 	}
// 	var rows, err = db.SelectSliceScan(sql, argMap)
// 	// log.Println(rows[0][0].(string))
// 	escola := &Escola{}
// 	errJson := json.Unmarshal([]byte(rows[0][0].(string)), escola)
// 	if errJson != nil {
// 		log.Println(errJson)
// 	}
// 	fmt.Printf("%+v\n", escola)
//
// 	for _, uni := range escola.Unidades {
// 		fmt.Printf("%+v\n", uni)
// 	}
// 	if err != nil {
// 		t.Error("Expected nothing, got ", err)
// 	}
//
// 	// return nil, err
// }

// var rows, err = h.DB.SelectSliceScan(sqlGetGrupoSacados, argMap)
// if err != nil {
//   if err == database.ErrNoRows {
//     return nil, ErrNoRowsGeneric
//   }
//   return nil, err
// }

func TestFluentd(t *testing.T) {
	logger, err := fluent.New(fluent.Config{FluentPort: 9880, FluentHost: "0.0.0.0"})
	if err != nil {
		fmt.Println(err)
	}
	defer logger.Close()
	tag := "myapp.access"
	var data = map[string]string{
		"foo":  "bar",
		"hoge": "hoge",
	}
	error := logger.Post(tag, data)
	// error := logger.PostWithTime(tag, time.Now(), data)
	if error != nil {
		t.Error("Expected nothing, got ", error)
	}
}
