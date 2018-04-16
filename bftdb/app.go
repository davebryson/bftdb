package bftdb

import (
	"fmt"

	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/tmlibs/merkle"
	"golang.org/x/crypto/ripemd160"
)

type Statement []byte

func (st Statement) Hash() []byte {
	hasher := ripemd160.New()
	hasher.Write(st)
	return hasher.Sum(nil)
}

func (st Statement) String() string {
	return string(st)
}

type App struct {
	abci.BaseApplication
	db        *DbWrapper
	entries   []merkle.Hasher
	lastBlock int64
	lastHash  []byte
}

func NewApp(dbo *DbWrapper) (*App, error) {
	// Create SQL connection
	return &App{
		db:        dbo,
		lastBlock: 0,
		lastHash:  []byte(""),
	}, nil
}

// Info something
func (app *App) Info(req abci.RequestInfo) abci.ResponseInfo {
	return abci.ResponseInfo{
		Data:             "bftsql",
		Version:          "1",
		LastBlockHeight:  app.lastBlock,
		LastBlockAppHash: app.lastHash,
	}
}

// InitChain called on startup
func (app *App) InitChain(req abci.RequestInitChain) abci.ResponseInitChain {
	table := `CREATE TABLE sample (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT)`
	e := app.db.Write(Statement(table))
	if e != nil {
		fmt.Printf("init %s\n", e.Error())
	}

	return abci.ResponseInitChain{}
}

// CheckTx validate incoming txs
func (app *App) CheckTx(tx []byte) abci.ResponseCheckTx {
	stmt := string(tx)
	if _, err := ValidateSql(stmt); err != nil {
		return abci.ResponseCheckTx{Code: 1, Log: "Bad SQL or DROP statement"}
	}

	return abci.ResponseCheckTx{Code: abci.CodeTypeOK}
}

// BeginBlock called on new block
func (app *App) BeginBlock(req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	// Setup entries for this block
	app.entries = make([]merkle.Hasher, 0)
	return abci.ResponseBeginBlock{}
}

// DeliverTx process txs
func (app *App) DeliverTx(tx []byte) abci.ResponseDeliverTx {
	statement := Statement(tx)
	e := app.db.Write(statement)
	if e != nil {
		return abci.ResponseDeliverTx{Code: 1}
	}

	app.entries = append(app.entries, statement)
	return abci.ResponseDeliverTx{Code: abci.CodeTypeOK}
}

// EndBlock called at end...
func (app *App) EndBlock(req abci.RequestEndBlock) abci.ResponseEndBlock {
	app.lastBlock = req.GetHeight()
	return abci.ResponseEndBlock{}
}

// Commit create hash of all txs
func (app *App) Commit() abci.ResponseCommit {
	if len(app.entries) > 0 {
		app.lastHash = merkle.SimpleHashFromHashers(app.entries)
	}
	app.entries = nil
	return abci.ResponseCommit{Data: app.lastHash}
}

// Query temp
func (app *App) Query(reqQuery abci.RequestQuery) abci.ResponseQuery {
	return abci.ResponseQuery{Log: "State Hash (base64 encoded)", Value: app.lastHash}
}
