package web

//  Disciplina ...
type Disciplina struct {
	IDDisciplina int      `json:"idDisciplina"`
	Nome         string   `json:"nome"`
	CargaHorario int      `json:"cargaHorario"`
	Ementa       string   `json:"ementa"`
	Descricao    string   `json:"descricao"`
	Serie        []*Serie `json:"series"`
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
