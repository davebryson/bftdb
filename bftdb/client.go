package bftdb

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	STATE_URL = "http://localhost:3000"
	QUERY_URL = "http://localhost:3000/query"
	OTHER_URL = "http://localhost:3000/stmt"
)

func encodeMsg(s string) ([]byte, error) {
	str := base64.StdEncoding.EncodeToString([]byte(s))
	jsonStr, e := json.Marshal([]string{str})
	if e != nil {
		return nil, e
	}
	return jsonStr, nil
}

func sendRequest(url string, method string, stmt []byte) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(stmt))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
	defer resp.Body.Close()

	fmt.Println("response Status :", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body   :", string(body))
}

func HandleSQL(stmt string) {
	stype, err := ValidateSql(stmt)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		jsonStr, e := encodeMsg(stmt)
		if e != nil {
			fmt.Println("Error encoding statement", e.Error())
		}
		switch stype {
		case SELECT:
			sendRequest(QUERY_URL, "POST", jsonStr)
		case OTHER:
			sendRequest(OTHER_URL, "POST", jsonStr)
		default:
			fmt.Println("Unknown request")
		}
	}
}

func HandleStatus() {
	sendRequest(STATE_URL, "GET", nil)
}
