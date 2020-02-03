package web

import (
	"log"

	"sab.io/escola-service/database"
)

// Unidade ...
type Unidade struct {
	IDUnidade int       `json:"idUnidade"`
	Nome      string    `json:"nome"`
	Endereco  *Endereco `json:"endereco"`
}

// CadastrarUnidade ...
func (u *Unidade) CadastrarUnidade(h *Handler, idEscola int, transDB *database.Transaction) error {
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

	log.Println("Endereco da unidade - ", u.Endereco.Logradouro)
	errEndereco := u.Endereco.CadastrarEndereco(h, transDB)
	if errEndereco != nil {
		erroTransDB = errEndereco
		return errEndereco
	}
	argMap := map[string]interface{}{
		"nome":        u.Nome,
		"id_escola":   idEscola,
		"id_endereco": u.Endereco.IDEndereco,
	}

	IDUnidade, errInsertUnidade := transDB.SelectSliceScan(database.SQLInsertUnidade, argMap)
	if errInsertUnidade != nil {
		erroTransDB = errInsertUnidade
		return errInsertUnidade
	}
	u.IDUnidade = rowNilInt(IDUnidade[0], 0)

	return nil
}
