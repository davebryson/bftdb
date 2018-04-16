package cmd

import (
	"context"
	"os"
	"time"

	"github.com/davebryson/bftdb/bftdb"
	cfg "github.com/tendermint/tendermint/config"
	nm "github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/proxy"
	"github.com/tendermint/tendermint/types"
	ttypes "github.com/tendermint/tendermint/types"
	cmn "github.com/tendermint/tmlibs/common"
	"github.com/tendermint/tmlibs/log"
)

const TESTDIR = "tmconfigs"

var (
	logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "main")
)

func doInit(config *cfg.Config, privValFile string) {
	if _, err := os.Stat(privValFile); os.IsNotExist(err) {
		// Validator file
		privValidator := ttypes.GenPrivValidatorFS(privValFile)
		//privValidator.SetFile(privValFile)
		privValidator.Save()

		// Genesis
		genFile := config.GenesisFile()
		// Default chain name
		chain_id := cmn.Fmt("bftdb-chain-%v", cmn.RandStr(6))

		// Create and save the genesis if it doesn't exist
		if _, err := os.Stat(genFile); os.IsNotExist(err) {
			// Set the chainid
			genDoc := ttypes.GenesisDoc{ChainID: chain_id}
			// Add the validators
			genDoc.Validators = []ttypes.GenesisValidator{ttypes.GenesisValidator{
				PubKey: privValidator.PubKey,
				Power:  10,
			}}
			genDoc.SaveAs(genFile)
		}
	}
}

func createTendermint(app *bftdb.App) *nm.Node {
	basedir := TESTDIR
	config := cfg.DefaultConfig()
	config.SetRoot(basedir)
	cfg.EnsureRoot(config.RootDir)

	privValFile := config.PrivValidatorFile()
	doInit(config, privValFile)

	privValidator := types.LoadOrGenPrivValidatorFS(privValFile)
	papp := proxy.NewLocalClientCreator(app)
	node, err := nm.NewNode(config, privValidator, papp,
		nm.DefaultGenesisDocProviderFunc(config),
		nm.DefaultDBProvider, logger)
	if err != nil {
		panic(err)
	}
	return node
}

func RunNode() {
	db, err := bftdb.NewDb()
	if err != nil {
		panic(err.Error())
	}

	logger.Info("Create db... starting server")
	app, err := bftdb.NewApp(db)
	if err != nil {
		panic(err.Error())
	}

	node := createTendermint(app)
	api := bftdb.NewQueryServer(db)

	node.Start()
	srv := api.Run()

	cmn.TrapSignal(func() {
		db.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("Server Shutdown:", err)
		}

		node.Stop()
	})
}
