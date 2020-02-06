package web

import (
	"encoding/json"
	"log"

	"sab.io/escola-service/database"
)

//  Disciplina ...
type Disciplina struct {
	IDDisciplina int       `json:"idDisciplina"`
	Nome         string    `json:"nome"`
	CargaHoraria int       `json:"cargaHoraria"`
	Descricao    string    `json:"descricao"`
	Ementas      []*Ementa `json:"ementas"`
}

type Ementa struct {
	Ementa   string `json:"ementa"`
	IDEmenta int    `json:"idEmenta"`
	Serie    *Serie `json:"serie"`
}
type Serie struct {
	IDSerie int    `json:"idSerie"`
	Tipo    string `json:"tipo"`
	Nome    string `json:"nome"`
}

// Matricula ...
type Matricula struct {
	ID      int     `json:"idMatricula"`
	IDTurma int     `json:"idTurma"`
	Serie   int     `json:"serie"`
	Status  string  `json:"status"`
	Nota    []*Nota `json:"nota"`
}

type Nota struct {
	ID           int        `json:"id"`
	IDDisciplina int        `json:"idDisciplina"`
	NotaFinal    int        `json:"notaFinal"`
	Bimestres    []Bimestre `json:"bimestres"`
}

type Bimestre struct {
	ID           int   `json:"id"`
	Nota         []int `json:"nota"`
	NotaBimestre int   `json:"notaBimestre"`
}

// GetDisciplinas ...
func GetDisciplinas(h *Handler, IDEscola int) ([]*Disciplina, error) {

	DB := h.DB
	argMap := map[string]interface{}{
		"id_escola": IDEscola,
	}
	rows, err := DB.SelectSliceScan(database.SQLGetDisciplinas, argMap)
	if err != nil {
		if err == database.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	disciplinas := make([]*Disciplina, len(rows))
	for i, row := range rows {
		d := rowNil(row, 0)
		log.Println("Disciplina Retornada -", d)
		disciplina := &Disciplina{}
		errJSON := json.Unmarshal([]byte(s), disciplina)
		if errJSON != nil {
			return nil, errJSON
		}
		disciplinas[i] = disciplina

	}
	return disciplinas, nil

}
