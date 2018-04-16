package bftdb

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/tendermint/tendermint/rpc/client"
)

// Routes
// /query  -> takes a query and returns result
// /stmt -> does the same
// / -> returns the lastes state
// PORT http addr
const PORT = ":3000"

type StatementRequest []string

func readStatementRequest(r io.Reader) (string, error) {
	var request StatementRequest
	data, e := ioutil.ReadAll(r)
	if e != nil {
		return "", e
	}

	e = json.Unmarshal(data, &request)
	if e != nil {
		return "", e
	}

	if len(request) < 0 {
		return "", e
	}

	// TODO: BASE64 DECODE the statement!
	p, e := base64.StdEncoding.DecodeString(request[0])
	if e != nil {
		return "", e
	}

	return string(p), nil
}

// QueryService base struct for QueryServer
type QueryService struct {
	db         *DbWrapper
	abciClient *client.HTTP //FIX with dep -> common problem with trace
}

func errorResp(msg string, writer io.Writer) {
	error := []string{msg}
	json.NewEncoder(writer).Encode(error)
}

// NewQueryServer create an http server that responds
// to SELECT statements
func NewQueryServer(db *DbWrapper) *QueryService {
	return &QueryService{
		db:         db,
		abciClient: client.NewHTTP("127.0.0.1:46657", "/websocket"),
	}
}

func (service *QueryService) QueryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		errorResp("Only POSTS", w)
		return
	}

	s, e := readStatementRequest(r.Body)
	if e != nil {
		errorResp(e.Error(), w)
		return
	}

	result, e := service.db.Read(s)
	if e != nil {
		errorResp(e.Error(), w)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func (service *QueryService) StatementHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		errorResp("Only POSTS", w)
		return
	}

	s, e := readStatementRequest(r.Body)
	if e != nil {
		errorResp(e.Error(), w)
		return
	}

	dtx, e := service.abciClient.BroadcastTxCommit([]byte(s))
	if e != nil {
		errorResp(e.Error(), w)
		return
	}

	json.NewEncoder(w).Encode(dtx)
}

func (service *QueryService) LatestState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	result, e := service.abciClient.ABCIQuery("", []byte(""))
	if e != nil {
		errorResp(e.Error(), w)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func (service *QueryService) Run() *http.Server {
	srv := &http.Server{Addr: PORT}

	router := http.NewServeMux()
	router.HandleFunc("/query", service.QueryHandler)
	router.HandleFunc("/stmt", service.StatementHandler)
	router.HandleFunc("/", service.LatestState)
	srv.Handler = router

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err.Error())
		}
	}()
	return srv
}
