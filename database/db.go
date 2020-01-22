//Package database é apenas
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
	//_ "github.com/marssis/gopgsqldriver"
	//driver postgres pure Go
	_ "github.com/lib/pq"
)

const (
	especFormatoSuportado              = "tdsf"
	tempoEsperaDefaultConnIdleDatabase = 30 //seconds
	tamDefaultPoolIdleConn             = 3
)

//ErrNoRows mensagem de erro quando não nenhum row
var ErrNoRows = errors.New("Não houve registros no retorno da query")

//constantes para o status da conexão do pool de conexões do database
const (
	statusPoolConnDisponivel = iota
	statusPoolConnEmUso
)

//DBLogger ...
type DBLogger interface {
	Trace(args ...interface{})
	Tracef(format string, args ...interface{})
	Notice(args ...interface{})
	Noticef(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Erro(args ...interface{})
	Errof(format string, args ...interface{})
}

//poolIdleConn representa uma conexão no pool
type poolIdleConn struct {
	db      *sql.DB
	timeout int64
	status  int
}

//OptionsDB contém o nome do banco de dados e informações de conexão.
type OptionsDB struct {
	DriverName        string
	IP                string
	Porta             int
	NomeDB            string
	User              string
	Senha             string
	Debug             bool //se true, ele loga os sqls
	Log               DBLogger
	TamPoolIdleConn   int
	TempoPoolIdleConn int
	Alias             string
	LogMinDuration    int
	MaxOpen           int
	DataSourceName    string
}

//DataBase é um extensao de sql.DB
type DataBase struct {
	*sql.DB
	options            *OptionsDB
	pool               []poolIdleConn
	canalSinalizaClose chan bool //sinaliza que o banco foi fechado a conexão
	muxPoolConn        sync.Mutex
}

//HostDB
func (this *DataBase) Host() string {
	return this.options.IP
}

//Porta
func (this *DataBase) Porta() int {
	return this.options.Porta
}

//Nome
func (this *DataBase) Nome() string {
	return this.options.NomeDB
}

//User
func (this *DataBase) User() string {
	return this.options.User
}

//Senha
func (this *DataBase) Senha() string {
	return this.options.Senha
}

//Open abre conexão com o banco de dados.
//Ela possui o mesmo comportamento de sql.Open, acrescido de um teste de conexão.
func (this *DataBase) Open() error {
	var (
		dataSourceName, traceMsg string
		erro                     error
	)

	if this.options.DriverName == "postgres" {
		//options='-c statement_timeout=100' - SET statement_timeout = '2s'
		if this.options.DataSourceName == "" {
			dataSourceName = fmt.Sprintf(`postgres://%s:%s@%s:%d/%s`,
				this.options.User, this.options.Senha, this.options.IP, this.options.Porta, this.options.NomeDB)
			fmt.Println("dataSourceName", "host=localhost port=5432 dbname=sabio user=postgres password=postgres sslmode=disable")
		} else {
			dataSourceName = this.options.DataSourceName
		}
	} else {
		return errors.New("Driver nao suportado:" + this.options.DriverName)
	}

	if this.options.TamPoolIdleConn <= 0 {
		this.options.TamPoolIdleConn = tamDefaultPoolIdleConn
	}

	if this.options.TempoPoolIdleConn <= 0 {
		this.options.TempoPoolIdleConn = tempoEsperaDefaultConnIdleDatabase
	}

	this.pool = make([]poolIdleConn, this.options.TamPoolIdleConn)
	this.canalSinalizaClose = make(chan bool)

	for i := 0; i < this.options.TamPoolIdleConn; i++ {
		// this.pool[i].db, erro = sql.Open(this.options.DriverName, dataSourceName)
		this.pool[i].db, erro = sql.Open("postgres", dataSourceName)
		if erro != nil {
			return erro
		}
		this.pool[i].status = statusPoolConnDisponivel
		this.pool[i].timeout = time.Now().Unix() + int64(this.options.TempoPoolIdleConn)
		this.pool[i].db.SetMaxIdleConns(1)

		traceMsg = fmt.Sprintf("Connecting(%d) to the database: %s:%s:%d...", i, this.options.Alias, this.options.IP, this.options.Porta)
		log.Println(traceMsg)
		if this.options.Log != nil {
			this.options.Log.Tracef(traceMsg)
		}
		log.Println("111111111")
		erro = this.pool[i].db.Ping()
		if erro != nil {
			return errors.New(dataSourceName + "-> " + erro.Error())
		}
		log.Println("2222222222")
		traceMsg = fmt.Sprintf("Connected(%d) to the database: %s:%s:%d [OK]", i, this.options.Alias, this.options.IP, this.options.Porta)
		if this.options.Log != nil {
			this.options.Log.Tracef(traceMsg)
		}
	}

	this.DB = this.pool[0].db
	if this.options.MaxOpen > 0 {
		this.DB.SetMaxOpenConns(this.options.MaxOpen)
	}
	this.initMonitorConnIdle()
	log.Println("333333333")
	return nil
}

//initMonitorConnIdle
func (this *DataBase) initMonitorConnIdle() {
	log.Println("ENTROU EM initMonitorConnIdle")
	go func() {
		for {
			select {
			case <-this.canalSinalizaClose: //qualquer coisa que chegar aqui, sair da goroutine
				log.Println(this.options.Alias + ":Banco fechado!")
				return
			case <-time.After(1 * time.Second): //chegou a hora de checar as conexões idle quanto a liberação
				this.liberarConnsIdle()
				break
			}
		}
	}()
}

//Close
func (this *DataBase) Close() {
	this.canalSinalizaClose <- true
	this.muxPoolConn.Lock()
	for i := 0; i < this.options.TamPoolIdleConn; i++ {
		this.pool[i].db.Close()
	}
	this.muxPoolConn.Unlock()
}

//liberarConnsIdle
func (this *DataBase) liberarConnsIdle() {
	this.muxPoolConn.Lock()
	defer this.muxPoolConn.Unlock()

	for i := 0; i < this.options.TamPoolIdleConn; i++ {
		if this.pool[i].timeout > 0 && (this.pool[i].status == statusPoolConnDisponivel) && (time.Now().Unix() > this.pool[i].timeout) {
			this.pool[i].db.SetMaxIdleConns(0)
			this.pool[i].timeout = 0
			//this.options.Log.Tracef("(%s) Liberado Conn Pool: %d", this.options.Alias, i)
			log.Printf("(%s) Liberado Conn Pool: %d", this.options.Alias, i)
		}
	}
}

//getConnPool pega um conexão disponivel no pool, se houver.
//Caso contrário, sempre a primeira do pool
func (this *DataBase) getConnPool() (int, *sql.DB) {
	this.muxPoolConn.Lock()
	defer this.muxPoolConn.Unlock()

	var i int
	for i = 0; i < this.options.TamPoolIdleConn; i++ {
		if this.pool[i].status == statusPoolConnDisponivel {
			this.pool[i].db.SetMaxIdleConns(1)
			this.pool[i].status = statusPoolConnEmUso
			return i, this.pool[i].db
		}
	}

	//como não houve nenhuma disponivel, logo this.pool[0].db.SetMaxIdleConns(1) já foi setado antes
	return -1, this.pool[0].db
}

//liberarConnPool
func (this *DataBase) liberarConnPool(indice int) {
	this.muxPoolConn.Lock()
	defer this.muxPoolConn.Unlock()

	if indice < 0 { //caso especial quando não houve nenhum disponível no pool
		//indice = 0
		return
	}
	this.pool[indice].status = statusPoolConnDisponivel
	this.pool[indice].timeout = time.Now().Unix() + int64(this.options.TempoPoolIdleConn)
}

//logarSql
func (this *DataBase) logarSql(query string, timeInicioConsulta time.Time, args ...interface{}) {
	if this.options.Debug {
		duracao := time.Now().Sub(timeInicioConsulta)
		if duracao > time.Duration(this.options.LogMinDuration)*time.Millisecond {
			logSql := fmt.Sprintf("[%s] [%s] SQL:\"%s\" Param($):%v", duracao.String(),
				formataMomento(timeInicioConsulta), query, args)
			if this.options.Log != nil {
				//this.options.Log.Trace(logSql)
				this.options.Log.Notice(logSql)
			}
		}
	}
}

//logarErroSql
func (this *DataBase) logarErroSql(query string, erro error, args ...interface{}) {
	//if this.options.Debug {
	logErroSql := fmt.Sprintf("Erro: %s == SQL:\"%s\" Param($):%v", erro.Error(), query, args)
	if this.options.Log != nil {
		this.options.Log.Erro(logErroSql)
	}
	//}
}

//StartTransaction inicia uma transação com o banco de dados
func (this *DataBase) StartTransaction() (*Transaction, error) {
	var i, db = this.getConnPool() //reservar até rollback ou commit

	tx, erro := db.Begin()
	if erro != nil {
		db.SetMaxIdleConns(0) //forçar desconexão
		this.liberarConnPool(i)
		return nil, erro
	}

	this.options.Log.Tracef("Begin")
	return &Transaction{db: this, tx: tx, closed: false, indicePool: i}, nil
}

//Query é um wrapper de sql.DB.Query
func (this *DataBase) Query(query string, args ...interface{}) (*sql.Rows, error) {
	var erroAux error

	var tempoInicio = time.Now()
	defer func(e error) {
		if e == nil {
			this.logarSql(query, tempoInicio, args)
		}
	}(erroAux)

	var i, db = this.getConnPool()
	defer this.liberarConnPool(i)
	log.Println("ANTES DA QUERY")
	rows, erro := db.Query(query, args...)
	log.Println(erro)
	if erro != nil {
		this.logarErroSql(query, erro, args)
		erroAux = erro
		db.SetMaxIdleConns(0) //forçar desconexão
		return nil, erro
	}
	log.Println("DEPOIS  DA QUERY")

	return rows, nil
}

//Query é um wrapper de sql.DB.QueryRow
func (this *DataBase) QueryRow(query string, args ...interface{}) *sql.Row {
	tempoInicio := time.Now()
	defer this.logarSql(query, tempoInicio, args)

	var i, db = this.getConnPool()
	defer this.liberarConnPool(i)

	row := db.QueryRow(query, args...)
	db.SetMaxIdleConns(0) //Como nao como verificar o erro aqui, melhor forçar desconexão no caso de ficar conexao remota inacessivel no caso de um kill no processo postgres
	return row
}

//Exec é um wrapper de sql.DB.Exec
func (this *DataBase) Exec(query string, args ...interface{}) (sql.Result, error) {
	var tempoInicio = time.Now()
	defer this.logarSql(query, tempoInicio, args)

	var i, db = this.getConnPool()
	defer this.liberarConnPool(i)

	result, erro := db.Exec(query, args...)
	if erro != nil {
		db.SetMaxIdleConns(0) //forçar desconexão
		this.logarErroSql(query, erro, args)
	}
	return result, erro
}

//Execute executa um query que não retornam rows
func (this *DataBase) QueryExec(query string, argMap map[string]interface{}, args ...interface{}) (sql.Result, error) {
	var sqlConsulta = montaQuery(query, argMap, args...)

	return this.Exec(sqlConsulta)
}

//Row representa uma tupla retornada por uma query como um slice
type Row []interface{}

func (this Row) Integer(coluna int) int {
	v, _ := strconv.Atoi(this[coluna].(string))
	return v
}

func (this Row) Float64(coluna int) float64 {
	v, _ := strconv.ParseFloat(this[coluna].(string), 64)
	return v
}

func (this Row) Bool(coluna int) bool {
	v, _ := strconv.ParseBool(this[coluna].(string))
	return v
}

func (this Row) String(coluna int) string {
	return this[coluna].(string)
}

//RowMap representa uma tupla retornada por uma query como um mapa
type RowMap map[string]interface{}

func (this RowMap) Integer(coluna string) int {
	v, _ := strconv.Atoi(this[coluna].(string))
	return v
}

//SelectSliceScan retorna as tuplas de uma query como um slice de []interface{}
//[]interface{} representa os campos da row
//argMap mapea cada chave em query para um valor do mapa. As chaves na query ficam entre '%('' e ')<especificador de formato, ESPEC_FORMATO_SUPORTADO>'
//args substitui os especificadores de formato normalmente, tais como %s, %d, %f
//Por exemplo:
//  qyery: select * from users where name = '%(name)s' and tipo = %d
//  argMap: map[string]interface{}{"name": "marcio"}
//  args: 1
//Observação:
//Quando o campo for null, haverá nil. Caso contrário, um tipo string
func (this *DataBase) SelectSliceScan(query string, argMap map[string]interface{}, args ...interface{}) ([]Row, error) {
	sqlConsulta := montaQuery(query, argMap, args...)
	log.Println(sqlConsulta)
	rows, erro := this.Query(sqlConsulta)
	if erro != nil {
		return nil, erro //em Query, já foi logado o erro
	}
	log.Println("RESULTADO RETORNADO - ", rows)
	defer rows.Close()

	tuplas, erro := SliceScan(rows, 1)
	if erro != nil {
		this.logarErroSql(sqlConsulta, erro, "SliceScan()")
		return nil, erro
	}

	if len(tuplas) <= 0 {
		return nil, ErrNoRows
	}

	return tuplas, nil
}

//SelectSliceScanCapacity ...
func (this *DataBase) SelectSliceScanCapacity(cap int, query string, argMap map[string]interface{}, args ...interface{}) ([]Row, error) {
	sqlConsulta := montaQuery(query, argMap, args...)

	rows, erro := this.Query(sqlConsulta)
	if erro != nil {
		return nil, erro //em Query, já foi logado o erro
	}
	defer rows.Close()

	tuplas, erro := SliceScan(rows, cap)
	if erro != nil {
		this.logarErroSql(sqlConsulta, erro, "SliceScan()")
		return nil, erro
	}

	if len(tuplas) <= 0 {
		return nil, ErrNoRows
	}

	return tuplas, nil
}

//SelectSliceScanObsoleto retorna as tuplas de uma query como um slice de []interface{}
//[]interface{} representa os campos da row
//argMap mapea cada chave em query para um valor do mapa. As chaves na query ficam entre '%('' e ')<especificador de formato, ESPEC_FORMATO_SUPORTADO>'
//args substitui os especificadores de formato normalmente, tais como %s, %d, %f
//Por exemplo:
//  qyery: select * from users where name = '%(name)s' and tipo = %d
//  argMap: map[string]interface{}{"name": "marcio"}
//  args: 1
//Observação:
//Quando o campo for null, haverá nil. Caso contrário, um tipo string
func (this *DataBase) SelectSliceScanObsoleto(query string, argMap map[string]interface{}, args ...interface{}) ([][]interface{}, error) {
	sqlConsulta := montaQuery(query, argMap, args...)

	rows, erro := this.Query(sqlConsulta)
	if erro != nil {
		return nil, erro //em Query, já foi logado o erro
	}
	defer rows.Close()

	tuplas, erro := SliceScanObsoleto(rows)
	if erro != nil {
		this.logarErroSql(sqlConsulta, erro, "SliceScan()")
		return nil, erro
	}

	return tuplas, nil
}

//SelectMapScan retorna as tuplas do comando select como um slice de mapa
//Funciona como SelectSliceScan, exceto pelo retorno
func (this *DataBase) SelectMapScan(query string, argMap map[string]interface{}, args ...interface{}) ([]RowMap, error) {
	sqlConsulta := montaQuery(query, argMap, args...)

	rows, erro := this.Query(sqlConsulta)
	if erro != nil {
		return nil, erro //em Query, já foi logado o erro
	}
	defer rows.Close()

	tuplas, erro := MapScan(rows)
	if erro != nil {
		this.logarErroSql(sqlConsulta, erro, "MapScan()")
		return nil, erro
	}

	if len(tuplas) <= 0 {
		return nil, ErrNoRows
	}

	return tuplas, nil
}

//SliceScan faz um scan em todas as rows, traduzindo-as em um simples slice de slice
func SliceScanObsoleto(rows *sql.Rows) ([][]interface{}, error) {
	colunas, erro := rows.Columns()
	if erro != nil {
		return nil, erro
	}

	tuplas := make([][]interface{}, 0)
	for rows.Next() {
		var i int
		tupla := make([]interface{}, len(colunas))
		for i = range tupla {
			tupla[i] = &sql.NullString{}
		}
		if erro = rows.Scan(tupla...); erro != nil {
			return nil, erro
		}

		//verificando se houve campo null
		for i = range colunas {
			ns := *(tupla[i].(*sql.NullString))
			if ns.Valid {
				tupla[i] = ns.String
			} else {
				tupla[i] = nil
			}
		}

		tuplas = append(tuplas, tupla)
	}

	if erro = rows.Err(); erro != nil {
		return nil, erro
	}

	return tuplas, nil
}

//SliceScan faz um scan em todas as rows, traduzindo-as em um simples slice de slice
func SliceScan(rows *sql.Rows, cap int) ([]Row, error) {
	colunas, erro := rows.Columns()
	if erro != nil {
		return nil, erro
	}

	if cap <= 0 {
		cap = 1
	}

	tuplas := make([]Row, 0, cap)
	for rows.Next() {
		var i int
		tupla := make(Row, len(colunas))
		for i = range tupla {
			tupla[i] = &sql.NullString{}
		}
		if erro = rows.Scan(tupla...); erro != nil {
			return nil, erro
		}

		//verificando se houve campo null
		for i = range colunas {
			ns := *(tupla[i].(*sql.NullString))
			if ns.Valid {
				tupla[i] = ns.String
			} else {
				tupla[i] = nil
			}
		}

		tuplas = append(tuplas, tupla)
	}

	if erro = rows.Err(); erro != nil {
		return nil, erro
	}

	return tuplas, nil
}

//MapScan faz um scan em todas as rows, traduzindo-as em um slice de mapa
func MapScan(rows *sql.Rows) ([]RowMap, error) {
	colunas, erro := rows.Columns()
	if erro != nil {
		return nil, erro
	}

	tuplas := make([]RowMap, 0, 1)
	for rows.Next() {
		tupla := make([]interface{}, len(colunas))
		for i := range tupla {
			tupla[i] = &sql.NullString{}
		}
		if erro = rows.Scan(tupla...); erro != nil {
			return nil, erro
		}

		//verificando se houve campo null
		tuplaMap := make(RowMap)
		for i, nomeColuna := range colunas {
			ns := *(tupla[i].(*sql.NullString))
			if ns.Valid {
				tuplaMap[nomeColuna] = ns.String
			} else {
				tuplaMap[nomeColuna] = nil
			}
		}

		tuplas = append(tuplas, tuplaMap)
	}

	if erro = rows.Err(); erro != nil {
		return nil, erro
	}

	return tuplas, nil
}

//NewDB ...
func NewDB(op *OptionsDB) *DataBase {
	db := &DataBase{}
	db.options = op
	if db.options.Alias == "" {
		db.options.Alias = db.options.NomeDB
	}

	return db
}

//montaQuery prepara uma query que possui especificadores de formato para usos em consultas, inserts e
//deletes
//query - string com os especificadores de formato
//mapa - é opcional. Utilizar nil quando não for utilizado. As chaves desse mapa aparecem entre % e o especificador
//       Por exemplo: %(nome)s
//args - substitui os especificadores de formato presentes na query
//
// Obserção: não conter espaços entre "%(" e ")" ao utilizar mapa na query
func montaQuery(query string, mapa map[string]interface{}, args ...interface{}) string {
	sql := query
	sql = strings.Replace(sql, "\n", "", -1)

	if mapa != nil {
		var (
			tokenParcial, token, especFormato string
			posicao                           int
		)
		for k, v := range mapa {
			tokenParcial = "%(" + k + ")"
			posicao = strings.Index(sql, tokenParcial)
			if posicao > 0 {
				//fmt.Println("sql[posicao:]", sql[posicao + len(tokenParcial):])
				especFormato = getEspecFormato(sql[posicao+len(tokenParcial):])
				//fmt.Println("especFormato", especFormato)
				if especFormato != "" {
					token = tokenParcial + especFormato
					//fmt.Println("token", token)
					especFormato = "%" + especFormato
					vlrMapa := fmt.Sprintf(especFormato, v)
					sql = strings.Replace(sql, token, vlrMapa, -1)
				}
			}
		}
	}

	//caso especial: tratar strings constantes com apóstrofo em seu conteúdo.
	var argsNew = make([]interface{}, 0, len(args))
	for _, arg := range args {
		var typ = reflect.TypeOf(arg)
		if typ.Kind() == reflect.String {
			var str = arg.(string)
			//str = strings.Trim(str, " ")
			if strings.HasPrefix(str, "'") && strings.HasSuffix(str, "'") {
				str = strings.Trim(str, "'")
				str = strings.Replace(str, "'", "''", -1)
				str = "'" + str + "'"
			} else {
				str = strings.Replace(str, "'", "''", -1)
			}
			argsNew = append(argsNew, str)
		} else {
			argsNew = append(argsNew, arg)
		}
	}

	//if len(args) > 0{
	sql = fmt.Sprintf(sql, argsNew...)
	sql = strings.Replace(sql, "(MISSING)", "", -1)
	sql = strings.Replace(sql, "!", "", 2)
	//}
	return sql
}

func getEspecFormato(s string) string {
	espec := make([]rune, 0, len(s))
	for _, ch := range s {
		if strings.ContainsAny(especFormatoSuportado, string(ch)) {
			espec = append(espec, ch)
			return string(espec)
		}
		espec = append(espec, ch)
	}

	return string(espec)
}

func formataMomento(t time.Time) string {
	momento := fmt.Sprintf("%02d/%02d/%04d %02d:%02d:%02d.%d",
		t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond()/1e6)
	return momento
}

//Transaction representa um transação no banco de dados
type Transaction struct {
	db         *DataBase
	tx         *sql.Tx
	closed     bool
	indicePool int
}

//Commit
func (this *Transaction) Commit() error {
	if !this.closed {
		this.db.liberarConnPool(this.indicePool)
		this.closed = true
		erro := this.tx.Commit()
		this.db.options.Log.Tracef("commit - %v", erro)
		return erro
	}

	return sql.ErrTxDone
}

//Rollback
func (this *Transaction) Rollback() error {
	if !this.closed {
		this.closed = true
		erro := this.tx.Rollback()
		this.db.options.Log.Tracef("rollback - %v", erro)
		this.db.SetMaxIdleConns(0) //forçar desconexão
		this.db.liberarConnPool(this.indicePool)
		return erro
	}

	return sql.ErrTxDone
}

//Exec é um wrapper de sql.Tx.Exec
func (this *Transaction) Exec(query string, argMap map[string]interface{}, args ...interface{}) (sql.Result, error) {
	var sqlConsulta = montaQuery(query, argMap, args...)

	var tempoInicio = time.Now()
	defer this.db.logarSql(sqlConsulta, tempoInicio, args)

	result, erro := this.tx.Exec(sqlConsulta)
	if erro != nil {
		this.db.logarErroSql(query, erro, args)
	}

	return result, erro
}

//Query é um wrapper de sql.Tx.Query
func (this *Transaction) Query(query string, argMap map[string]interface{}, args ...interface{}) (*sql.Rows, error) {
	var sqlConsulta = montaQuery(query, argMap, args...)

	var tempoInicio = time.Now()
	defer this.db.logarSql(sqlConsulta, tempoInicio, args)

	rows, erro := this.tx.Query(sqlConsulta)
	if erro != nil {
		this.db.logarErroSql(query, erro, args)
		return nil, erro
	}

	return rows, nil
}

//SelectSliceScan
func (this *Transaction) SelectSliceScan(query string, argMap map[string]interface{}, args ...interface{}) ([]Row, error) {
	var sqlConsulta = montaQuery(query, argMap, args...)
	log.Println("SQLConsutal - ", sqlConsulta)
	var tempoInicio = time.Now()
	defer this.db.logarSql(sqlConsulta, tempoInicio, args)

	rows, erro := this.tx.Query(sqlConsulta)
	if erro != nil {
		this.db.logarErroSql(query, erro, args)
		return nil, erro
	}
	defer rows.Close()

	tuplas, erro := SliceScan(rows, 1)
	if erro != nil {
		this.db.logarErroSql(sqlConsulta, erro, "Transaction->SliceScan()")
		return nil, erro
	}

	if len(tuplas) <= 0 {
		return nil, ErrNoRows
	}

	return tuplas, nil
}
