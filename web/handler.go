package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"sab.io/escola-service/cache"
	"sab.io/escola-service/database"
	"sab.io/escola-service/logger"
	"sab.io/escola-service/messaging"
	"sab.io/escola-service/storage"
)

var (
	timeDefaultCache = time.Minute * 10
)

// Handler ...
type Handler struct {
	// Relic        newrelic.Application
	// ClientRedis  *redis.Client
	DB *database.DataBase
	// HostsVálidos []string
	// Log          *snetlog.Log
	// GodinAuthURL string
	// Upgrader     websocket.Upgrader
	// EmailConf    map[string]string
	Message *messaging.IMessageClient
	Cache   *cache.ICacheClient
	Logger  *logger.ILogger
	Storage *storage.IStorage
}

func (h *Handler) Test(m *messaging.MessageParam) *messaging.MessageParam {
	log.Println("Mensagem recebida")
	return m
}

// CadastrarEscola ...
func (h *Handler) CadastrarEscola(m *messaging.MessageParam) *messaging.MessageParam {
	logg := *h.Logger
	memCache := *h.Cache
	escola := &Escola{}
	errJSON := json.Unmarshal(m.Body, escola)
	if errJSON != nil {
		log.Println(errJSON)
		logg.Send(logger.ERROR, errJSON.Error(), m.IDOperation)
		return nil
	}
	errCad := escola.CadastrarEscola(h, nil)
	if errCad != nil {
		logg.Send(logger.ERROR, errCad.Error(), m.IDOperation)
		m.Type = messaging.TYPE_ERROR
		m.Body = []byte(errCad.Error())
		return m
	}
	errCache := memCache.SetValue(fmt.Sprintf("escola:%d", escola.IDEscola), escola, timeDefaultCache)
	if errCache != nil {
		logg.Send(logger.WARNING, errCache.Error(), m.IDOperation)
	}
	body, errJSONB := json.Marshal(escola)
	if errJSONB != nil {
		log.Println(errJSONB)
		logg.Send(logger.ERROR, errJSONB.Error(), m.IDOperation)
		return nil
	}
	m.Type = messaging.TYPE_RESPONSE
	m.Body = body

	return m
}

// GetEscolas ...
func (h *Handler) GetEscolas(m *messaging.MessageParam) *messaging.MessageParam {
	logg := *h.Logger
	escolas, errEscolas := GetEscolas(h)
	if errEscolas != nil {
		logg.Send(logger.ERROR, errEscolas.Error(), m.IDOperation)
		m.Type = messaging.TYPE_ERROR
		m.Body = []byte(errEscolas.Error())
	}
	body, errJSONB := json.Marshal(escolas)
	if errJSONB != nil {
		log.Println(errJSONB)
		logg.Send(logger.ERROR, errJSONB.Error(), m.IDOperation)
		return nil
	}
	m.Type = messaging.TYPE_RESPONSE
	m.Body = body
	return m
}

// GetEscola ...
func (h *Handler) GetEscola(m *messaging.MessageParam) *messaging.MessageParam {
	logg := *h.Logger
	memCache := *h.Cache
	IDEscola := m.Params["idEscola"]
	escola := &Escola{}
	errCache := memCache.GetValue(fmt.Sprintf("escola:%d", IDEscola), escola)
	if errCache != nil {
		logg.Send(logger.WARNING, errCache.Error(), m.IDOperation)
	}
	if escola.IDEscola == 0 {
		escolaB, errEscola := GetEscola(h, IDEscola)
		if errEscola != nil {
			logg.Send(logger.ERROR, errEscola.Error(), m.IDOperation)
			m.Type = messaging.TYPE_ERROR
			m.Body = []byte(errEscola.Error())
			return m
		}
		escola = escolaB

	}

	errCache = memCache.SetValue(fmt.Sprintf("escola:%d", IDEscola), escola, timeDefaultCache)
	if errCache != nil {
		logg.Send(logger.WARNING, errCache.Error(), m.IDOperation)
	}

	body, errJSONB := json.Marshal(escola)
	if errJSONB != nil {
		log.Println(errJSONB)
		logg.Send(logger.ERROR, errJSONB.Error(), m.IDOperation)
		return nil
	}

	m.Type = messaging.TYPE_RESPONSE
	m.Body = body
	return m
}

func (h *Handler) CadastrarUnidade(m *messaging.MessageParam) *messaging.MessageParam {
	logg := *h.Logger
	memCache := *h.Cache
	unidade := &Unidade{}
	idEscola := m.Params["idEscola"]
	errJSON := json.Unmarshal(m.Body, unidade)
	if errJSON != nil {
		log.Println(errJSON)
		logg.Send(logger.ERROR, errJSON.Error(), m.IDOperation)
		return nil
	}
	errCad := unidade.CadastrarUnidade(h, idEscola, nil)
	if errCad != nil {
		logg.Send(logger.ERROR, errCad.Error(), m.IDOperation)
		m.Type = messaging.TYPE_ERROR
		m.Body = []byte(errCad.Error())
		return m
	}
	escola := &Escola{}
	errCache := memCache.GetValue(fmt.Sprintf("escola:%d", idEscola), escola)
	if errCache != nil {
		logg.Send(logger.WARNING, errCache.Error(), m.IDOperation)
	}
	if escola.IDEscola != 0 {
		escola.Unidades = append(escola.Unidades, unidade)
		errCache := memCache.SetValue(fmt.Sprintf("escola:%d", idEscola), escola, timeDefaultCache)
		if errCache != nil {
			logg.Send(logger.WARNING, errCache.Error(), m.IDOperation)
		}
	}
	body, errJSONB := json.Marshal(unidade)
	if errJSONB != nil {
		log.Println(errJSONB)
		logg.Send(logger.ERROR, errJSONB.Error(), m.IDOperation)
		return nil
	}

	m.Type = messaging.TYPE_RESPONSE
	m.Body = body
	return m
}

// CadastrarDisciplina ...
func (h *Handler) CadastrarDisciplina(m *messaging.MessageParam) *messaging.MessageParam {
	logg := *h.Logger
	memCache := *h.Cache
	disciplina := &Disciplina{}
	idEscola := m.Params["idEscola"]
	errJSON := json.Unmarshal(m.Body, disciplina)
	if errJSON != nil {
		log.Println(errJSON)
		logg.Send(logger.ERROR, errJSON.Error(), m.IDOperation)
		return nil
	}
	errCad := disciplina.CadastrarDisciplina(h, idEscola, nil)
	if errCad != nil {
		logg.Send(logger.ERROR, errCad.Error(), m.IDOperation)
		m.Type = messaging.TYPE_ERROR
		m.Body = []byte(errCad.Error())
		return m
	}
	errCache := memCache.SetValue(fmt.Sprintf("disciplina:%d", disciplina.IDDisciplina), disciplina, timeDefaultCache)
	if errCache != nil {
		logg.Send(logger.WARNING, errCache.Error(), m.IDOperation)
	}

	body, errJSONB := json.Marshal(disciplina)
	if errJSONB != nil {
		log.Println(errJSONB)
		logg.Send(logger.ERROR, errJSONB.Error(), m.IDOperation)
		return nil
	}
	m.Type = messaging.TYPE_RESPONSE
	m.Body = body

	return m
}

// GetDisciplinas ...
func (h *Handler) GetDisciplinas(m *messaging.MessageParam) *messaging.MessageParam {
	logg := *h.Logger
	idEscola := m.Params["idEscola"]
	disciplinas, errDisciplinas := GetDisciplinas(h, idEscola)
	if errDisciplinas != nil {
		logg.Send(logger.ERROR, errDisciplinas.Error(), m.IDOperation)
		m.Type = messaging.TYPE_ERROR
		m.Body = []byte(errDisciplinas.Error())
	}
	body, errJSONB := json.Marshal(disciplinas)
	if errJSONB != nil {
		log.Println(errJSONB)
		logg.Send(logger.ERROR, errJSONB.Error(), m.IDOperation)
		return nil
	}
	m.Type = messaging.TYPE_RESPONSE
	m.Body = body
	return m
}

// GetDisciplina ...
func (h *Handler) GetDisciplina(m *messaging.MessageParam) *messaging.MessageParam {
	logg := *h.Logger
	memCache := *h.Cache
	IDEscola := m.Params["idEscola"]
	IDDisciplina := m.Params["idDisciplina"]
	disciplina := &Disciplina{}
	errCache := memCache.GetValue(fmt.Sprintf("disciplina:%d", IDDisciplina), disciplina)
	if errCache != nil {
		logg.Send(logger.WARNING, errCache.Error(), m.IDOperation)
	}
	if disciplina.IDDisciplina == 0 {
		disciplinaB, errDisciplina := GetDisciplinaByID(h, IDEscola, IDDisciplina)
		if errDisciplina != nil {
			logg.Send(logger.ERROR, errDisciplina.Error(), m.IDOperation)
			m.Type = messaging.TYPE_ERROR
			m.Body = []byte(errDisciplina.Error())
			return m
		}
		disciplina = disciplinaB

	}
	body, errJSONB := json.Marshal(disciplina)
	if errJSONB != nil {
		log.Println(errJSONB)
		logg.Send(logger.ERROR, errJSONB.Error(), m.IDOperation)
		return nil
	}

	m.Type = messaging.TYPE_RESPONSE
	m.Body = body
	return m
}

// AtualizaDisciplina ...
func (h *Handler) AtualizaDisciplina(m *messaging.MessageParam) *messaging.MessageParam {
	logg := *h.Logger
	memCache := *h.Cache
	disciplina := &Disciplina{}
	idEscola := m.Params["idEscola"]
	errJSON := json.Unmarshal(m.Body, disciplina)
	if errJSON != nil {
		log.Println(errJSON)
		logg.Send(logger.ERROR, errJSON.Error(), m.IDOperation)
		return nil
	}
	errCad := disciplina.AtualizarDisciplina(h, idEscola, nil)
	if errCad != nil {
		logg.Send(logger.ERROR, errCad.Error(), m.IDOperation)
		m.Type = messaging.TYPE_ERROR
		m.Body = []byte(errCad.Error())
		return m
	}
	errCache := memCache.SetValue(fmt.Sprintf("disciplina:%d", disciplina.IDDisciplina), disciplina, timeDefaultCache)
	if errCache != nil {
		logg.Send(logger.WARNING, errCache.Error(), m.IDOperation)
	}

	body, errJSONB := json.Marshal(disciplina)
	if errJSONB != nil {
		log.Println(errJSONB)
		logg.Send(logger.ERROR, errJSONB.Error(), m.IDOperation)
		return nil
	}
	m.Type = messaging.TYPE_RESPONSE
	m.Body = body

	return m
}

// DeleteDisciplina ...
func (h *Handler) DeleteDisciplina(m *messaging.MessageParam) *messaging.MessageParam {
	logg := *h.Logger
	memCache := *h.Cache
	disciplina := &Disciplina{}
	idEscola := m.Params["idEscola"]
	idDisciplina := m.Params["idDisciplina"]

	err := memCache.GetValue(fmt.Sprintf("disciplina:%d", idDisciplina), disciplina)
	if err != nil {
		logg.Send(logger.WARNING, err.Error(), m.IDOperation)
	}
	disciplina, err = GetDisciplinaByID(h, idEscola, idDisciplina)
	if err != nil {
		logg.Send(logger.ERROR, err.Error(), m.IDOperation)
	}

	// errJSON := json.Unmarshal(m.Body, disciplina)
	// if errJSON != nil {
	// 	log.Println(errJSON)
	// 	logg.Send(logger.ERROR, errJSON.Error(), m.IDOperation)
	// 	return nil
	// }
	err = disciplina.DeletarDisciplina(h, idEscola, nil)
	if err != nil {
		logg.Send(logger.ERROR, err.Error(), m.IDOperation)
		m.Type = messaging.TYPE_ERROR
		m.Info = err.Error()
		m.Status = http.StatusNotFound
		m.Body = []byte("")
		return m
	}
	err = memCache.DelValue(fmt.Sprintf("disciplina:%d", disciplina.IDDisciplina))
	if err != nil {
		logg.Send(logger.WARNING, err.Error(), m.IDOperation)
	}

	m.Status = http.StatusOK
	m.Type = messaging.TYPE_RESPONSE
	m.Info = DisciplinaDel
	m.Body = []byte("")

	return m
}

func (h *Handler) DeleteEmenta(m *messaging.MessageParam) *messaging.MessageParam {
	logg := *h.Logger
	memCache := *h.Cache
	ementa := &Ementa{}
	idDisciplina := m.Params["idDisciplina"]
	errJSON := json.Unmarshal(m.Body, ementa)
	if errJSON != nil {
		log.Println(errJSON)
		logg.Send(logger.ERROR, errJSON.Error(), m.IDOperation)
		return nil
	}
	errCad := ementa.DeleteEmenta(h, idDisciplina, nil)
	if errCad != nil {
		logg.Send(logger.ERROR, errCad.Error(), m.IDOperation)
		m.Type = messaging.TYPE_ERROR
		m.Body = []byte(errCad.Error())
		return m
	}
	disciplina := &Disciplina{}
	errCache := memCache.GetValue(fmt.Sprintf("disciplina:%d", idDisciplina), disciplina)
	if errCache != nil {
		logg.Send(logger.WARNING, errCache.Error(), m.IDOperation)
	}
	if disciplina.IDDisciplina != 0 {
		for i, ement := range disciplina.Ementas {
			if ement.IDEmenta == ementa.IDEmenta {
				disciplina.Ementas[i] = ementa
			}
		}
		errCache := memCache.SetValue(fmt.Sprintf("disciplina:%d", idDisciplina), disciplina, timeDefaultCache)
		if errCache != nil {
			logg.Send(logger.WARNING, errCache.Error(), m.IDOperation)
		}
	}
	body, errJSONB := json.Marshal(ementa)
	if errJSONB != nil {
		log.Println(errJSONB)
		logg.Send(logger.ERROR, errJSONB.Error(), m.IDOperation)
		return nil
	}

	m.Type = messaging.TYPE_RESPONSE
	m.Body = body
	return m
}
