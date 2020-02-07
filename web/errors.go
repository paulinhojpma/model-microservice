package web

import (
	"errors"
)

var (
	ErrorNoUnidade = errors.New("Nenhuma unidade enviada para cadastro da escola")
	ErrorAtualizar = errors.New("Não foi possível atualizar o recurso")
)
