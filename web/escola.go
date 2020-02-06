package web

import (
	"encoding/json"
	"log"

	"sab.io/escola-service/database"
)

// Escola ...
type Escola struct {
	IDEscola int        `json:"idEscola"`
	Nome     string     `json:"nome"`
	Unidades []*Unidade `json:"unidades"`
	Cnpj     string     `json:"cnpj"`
}

type Professor struct {
	ID   int    `json:"idProfessor"`
	Nome string `json:"nome"`
	Cpf  string `json:"cpf"`
}

type AnoLetivo struct {
	ID  int `json:"idAnoLetivo"`
	Ano int `json:"ano"`
}

// GetEscolas ...
func GetEscolas(h *Handler) ([]*Escola, error) {
	escolas := make([]*Escola, 0)
	DB := h.DB
	rows, err := DB.SelectSliceScan(database.SQLGetEscolas, nil)
	if err != nil {
		if err == database.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	for _, row := range rows {
		s := rowNil(row, 0)
		log.Println("ESCOLA Retornada -", s)
		escola := &Escola{}
		errJSON := json.Unmarshal([]byte(s), escola)
		if errJSON != nil {
			return nil, errJSON
		}
		escolas = append(escolas, escola)

	}
	return escolas, nil

}

// CadastrarEscola ...
func (e *Escola) CadastrarEscola(h *Handler, transDB *database.Transaction) error {
	DB := h.DB
	var (
		erroTransDB error
	)
	if transDB == nil {
		if transDB, erroTransDB = DB.StartTransaction(); erroTransDB != nil {
			return nil
		}
		defer func() {
			if transDB != nil {
				if erroTransDB != nil {
					transDB.Rollback()
				} else {
					transDB.Commit()
				}
			}
		}()

	}

	argMap := map[string]interface{}{
		"nome": e.Nome,
		"cnpj": e.Cnpj,
	}

	IDEscola, errInsertEscola := transDB.SelectSliceScan(database.SQLInsertEscola, argMap)
	if errInsertEscola != nil {
		erroTransDB = errInsertEscola
		return errInsertEscola
	}
	e.IDEscola = rowNilInt(IDEscola[0], 0)
	if len(e.Unidades) == 0 {
		erroTransDB = ErrorNoUnidade
		return ErrorNoUnidade
	}
	for _, u := range e.Unidades {
		errUn := u.CadastrarUnidade(h, e.IDEscola, transDB)
		if errUn != nil {
			erroTransDB = errUn
			return errUn
		}
	}

	return nil
}

// GetEscola ...
func GetEscola(h *Handler, idEscola int) (*Escola, error) {
	escola := &Escola{}
	DB := h.DB
	argMap := map[string]interface{}{
		"id_escola": idEscola,
	}
	row, err := DB.SelectSliceScan(database.SQLGetEscolaByID, argMap)
	if err != nil {
		if err == database.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	s := rowNil(row[0], 0)
	log.Println(s)
	errJSON := json.Unmarshal([]byte(s), escola)
	if errJSON != nil {
		return nil, errJSON
	}

	return escola, nil
}
