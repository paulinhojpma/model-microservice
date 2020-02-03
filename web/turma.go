package web

// Turma ...
type Turma struct {
	ID           int        `json:"id"`
	AnoEscolar   int        `json:"anoEscolar"`
	Status       string     `json:"status"`
	Serie        int        `json:"serie"`
	Horarios     []*Horario `json:"horarios"`
	IDExercicios []int      `json:"idExercicios"`
}

// Horario ...
type Horario struct {
	IDHorario    int
	HoraInicio   int   `json:"horaInicio"`
	MinutoInicio int   `json:"minutoInicio"`
	Sala         *Sala `json:"sala"`
	IDDisciplina int   `json:"idDisciplina"`
	IDProfessor  int   `json:"idProfessor"`
}

// Sala ...
type Sala struct {
	ID         int    `json:"id"`
	Nome       string `json:"nome"`
	Status     string `json:"status"`
	Capacidade int    `json:"capacidade"`
}
