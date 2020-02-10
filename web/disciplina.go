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
	Ementa       string `json:"ementa"`
	IDEmenta     int    `json:"idEmenta"`
	Ativo        bool   `json:"ativo"`
	CargaHoraria int    `json:"cargaHoraria"`
	Serie        *Serie `json:"serie"`
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
		errJSON := json.Unmarshal([]byte(d), disciplina)
		if errJSON != nil {
			return nil, errJSON
		}
		disciplinas[i] = disciplina

	}
	return disciplinas, nil

}

// GetDisciplinas ...
func GetDisciplinaByID(h *Handler, IDEscola, IDDisciplina int) (*Disciplina, error) {

	DB := h.DB
	argMap := map[string]interface{}{
		"id_escola":     IDEscola,
		"id_disciplina": IDDisciplina,
	}
	rows, err := DB.SelectSliceScan(database.SQLGetDisciplinaByID, argMap)
	if err != nil {
		if err == database.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	d := rowNil(rows[0], 0)
	log.Println("Disciplina Retornada -", d)
	disciplina := &Disciplina{}
	errJSON := json.Unmarshal([]byte(d), disciplina)
	if errJSON != nil {
		return nil, errJSON
	}

	return disciplina, nil
}

// CadastrarDisciplina ...
func (d *Disciplina) CadastrarDisciplina(h *Handler, idEscola int, transDB *database.Transaction) error {
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
		"nome":      d.Nome,
		"descricao": d.Descricao,
		"id_escola": idEscola,
	}

	rowIdDisciplina, errInsertDisciplina := transDB.SelectSliceScan(database.SQLInsertDisciplina, argMap)
	if errInsertDisciplina != nil {
		erroTransDB = errInsertDisciplina
		return errInsertDisciplina
	}

	d.IDDisciplina = rowNilInt(rowIdDisciplina[0], 0)
	if len(d.Ementas) != 0 {
		for _, e := range d.Ementas {
			errEm := e.cadastrarEmenta(h, d.IDDisciplina, transDB)
			if errEm != nil {
				erroTransDB = errEm
				return errEm
			}
		}
	}

	return nil
}

// AtualizarDisciplina ...
func (d *Disciplina) AtualizarDisciplina(h *Handler, idEscola int, transDB *database.Transaction) error {
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
		"nome":          d.Nome,
		"descricao":     d.Descricao,
		"id_escola":     idEscola,
		"id_disciplina": d.IDDisciplina,
	}

	rowIdDisciplina, errInsertDisciplina := transDB.SelectSliceScan(database.SQLUpdateDisciplina, argMap)
	if errInsertDisciplina != nil {
		erroTransDB = errInsertDisciplina
		return errInsertDisciplina
	}

	idDisciplina := rowNilInt(rowIdDisciplina[0], 0)
	if idDisciplina == 0 {
		erroTransDB = ErrorAtualizar
		return ErrorAtualizar
	}
	if len(d.Ementas) != 0 {
		for _, e := range d.Ementas {
			errEm := e.cadastrarEmenta(h, d.IDDisciplina, transDB)
			if errEm != nil {
				erroTransDB = errEm
				return errEm
			}
		}
	}

	return nil
}

// DeletarDisciplina
func (d *Disciplina) DeletarDisciplina(h *Handler, idEscola int, transDB *database.Transaction) error {
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
		"id_escola":     idEscola,
		"id_disciplina": d.IDDisciplina,
	}

	rowIdDisciplina, errInsertDisciplina := transDB.SelectSliceScan(database.SQLDeleteDisciplina, argMap)
	if errInsertDisciplina != nil {
		erroTransDB = errInsertDisciplina
		return errInsertDisciplina
	}

	idDisciplina := rowNilInt(rowIdDisciplina[0], 0)
	if idDisciplina == 0 {
		erroTransDB = ErrorAtualizar
		return ErrorAtualizar
	}
	if len(d.Ementas) != 0 {
		for _, e := range d.Ementas {
			errEm := e.DeleteEmenta(h, d.IDDisciplina, transDB)
			if errEm != nil {
				erroTransDB = errEm
				return errEm
			}
		}
	}

	return nil
}

func (e *Ementa) cadastrarEmenta(h *Handler, idDisciplina int, transDB *database.Transaction) error {
	argMap := map[string]interface{}{
		"carga_horaria": e.CargaHoraria,
		"ementa":        e.Ementa,
		"id_serie":      e.Serie.IDSerie,
		"id_disciplina": idDisciplina,
		"ativo":         e.Ativo,
	}

	rowIdEmenta, errInsertEmenta := transDB.SelectSliceScan(database.SQLInsertEmenta, argMap)
	if errInsertEmenta != nil {
		return errInsertEmenta
	}
	e.IDEmenta = rowNilInt(rowIdEmenta[0], 0)
	return nil
}

func (e *Ementa) DeleteEmenta(h *Handler, idDisciplina int, transDB *database.Transaction) error {
	argMap := map[string]interface{}{
		"id_disciplina": idDisciplina,
	}

	rowIdEmenta, errInsertEmenta := transDB.SelectSliceScan(database.SQLDeleteEmenta, argMap)
	if errInsertEmenta != nil {
		return errInsertEmenta
	}
	e.IDEmenta = rowNilInt(rowIdEmenta[0], 0)
	return nil
}
