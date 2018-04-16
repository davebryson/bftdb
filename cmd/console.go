package cmd

import (
	"fmt"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/davebryson/bftdb/bftdb"
)

const HELP = `Commands:
- help
- state :  Show the lastest state on the blockchain

Otherwise, simply enter SQL commands at the prompt, e.g. 'select * from users'`

const BANNER = `
 ___ ___ _____    ___  ___  _    ___ _____ ___
| _ ) __|_   _|__/ __|/ _ \| |  |_ _|_   _| __|
| _ \ _|  | ||___\__ \ (_) | |__ | |  | | | _|   Tendermint + SQLite3
|___/_|   |_|    |___/\__\_\____|___| |_| |___|
Press ctrl-D to exit
Type 'help' for more info`

func executor(in string) {
	cmd := strings.TrimSpace(in)
	switch cmd {
	case "help", "h":
		fmt.Println(HELP)
	case "state", "s":
		bftdb.HandleStatus()
	default:
		bftdb.HandleSQL(cmd)
	}
}

func completer(in prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

func RunConsole() {
	p := prompt.New(executor,
		completer,
		prompt.OptionTitle("bft-sqlite-prompt"),
		prompt.OptionHistory([]string{}),
		prompt.OptionPrefixTextColor(prompt.Yellow),
		prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionBGColor(prompt.DarkGray))

	fmt.Println(BANNER)
	p.Run()
}
