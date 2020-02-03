package web

import "sab.io/escola-service/database"

// Endereco ...
type Endereco struct {
	IDEndereco  int    `json:"idEndereco"`
	Logradouro  string `json:"logradouro"`
	Numero      string `json:"numero"`
	Cep         string `json:"cep"`
	Bairro      string `json:"bairro"`
	Complemento string `json:"complemento"`
	UF          string `json:"uf"`
	Cidade      string `json:"cidade"`
}

// CadastrarEndereco ...
func (e *Endereco) CadastrarEndereco(h *Handler, transDB *database.Transaction) error {

	argMap := map[string]interface{}{
		"logradouro":  e.Logradouro,
		"numero":      e.Numero,
		"cep":         e.Cep,
		"bairro":      e.Bairro,
		"complemento": e.Complemento,
		"uf":          e.UF,
		"cidade":      e.Cidade,
	}
	rowIDEndereco, errInsertEndereco := transDB.SelectSliceScan(database.SQLInsertEndereco, argMap)
	if errInsertEndereco != nil {
		return errInsertEndereco
	}
	e.IDEndereco = rowNilInt(rowIDEndereco[0], 0)

	return nil
}
