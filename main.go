package main

import (
	"encoding/json"
	"fmt"
	"github.com/rockjoon/nomadcoin/blockchain"
	"net/http"
)

const port string = ":4000"

type URLDescription struct {
	URL         URL    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

type URL string

type AddBlockBody struct {
	Data string
}

func (u URL) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

func main() {
	//explorer.Start()
	http.HandleFunc("/", documentation)
	http.HandleFunc("/blocks", blocks)
	http.ListenAndServe(port, nil)
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rw.Header().Add("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(blockchain.GetBlockChain().AllBlocks())
	case http.MethodPost:
		var addBlockBody AddBlockBody
		json.NewDecoder(r.Body).Decode(&addBlockBody)
		blockchain.GetBlockChain().AddBlock(addBlockBody.Data)
		rw.WriteHeader(http.StatusCreated)
	}
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []URLDescription{
		{
			"/", "GET", "See Documentation", "",
		},
		{"/blocks", "POST", "Add a block", "data:string"},
	}
	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(data)
}
