package storage

import (
	"errors"
)

var (
	// ErrTipoArquivoInvalido ...
	ErrTipoArquivoInvalido = errors.New("Tipo de arquivo inválido")
	ErrParametroArquivo    = errors.New("Parametros de nome de arquivo invalido")
	ErrMaxFileSize         = errors.New("Arquivo enviado muito grande")
)
