package web

// Aluno ...
type Aluno struct {
	ID             int      `json:"idAluno"`
	Nome           string   `json:"nome"`
	Cpf            string   `json:"cpf"`
	ValMatricula      string   `json:"valMatricula"`
	DataNascimento string   `json:"dataNascimento"`
	TurmaID        int      `json:"turmaID"`
	Endereco       Endereco `json:"endereco"`
	Matricula []*Matricula `json:"matriculas"`

}
