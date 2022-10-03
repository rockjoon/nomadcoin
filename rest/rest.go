package rest

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rockjoon/nomadcoin/blockchain"
	"github.com/rockjoon/nomadcoin/utils"
	"log"
	"net/http"
)

var port string

type URL string

type urlDescription struct {
	Url         URL    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

type addBlockBody struct {
	Data string
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

func (u URL) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func Start(aPort int) {
	port = fmt.Sprintf(":%d", aPort)
	router := mux.NewRouter()
	router.Use(jsonContentTypeMiddleware)
	router.HandleFunc("/", documentation).Methods(http.MethodGet)
	router.HandleFunc("/blocks", blocks).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods(http.MethodGet)
	log.Printf("listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	block, err := blockchain.FindBlock(hash)
	encoder := json.NewEncoder(rw)
	if err != nil {
		encoder.Encode(errorResponse{fmt.Sprint(err)})
	} else {
		encoder.Encode(block)
	}

}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(rw).Encode(blockchain.GetBlockChain().AllBlocks())
	case http.MethodPost:
		var addBlockBody addBlockBody
		utils.HandleError(json.NewDecoder(r.Body).Decode(&addBlockBody))
		blockchain.GetBlockChain().AddBlock(addBlockBody.Data)
		rw.WriteHeader(http.StatusCreated)
	}
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{
			Url:         URL("/"),
			Method:      "GET",
			Description: "See documentation",
		},
		{
			Url:         URL("/blocks"),
			Method:      "GET",
			Description: "See all blocks",
		},
		{
			Url:         URL("/blocks"),
			Method:      "POST",
			Description: "Add a block",
			Payload:     "Block data",
		},
		{
			Url:         URL("/blocks/{hash}"),
			Method:      "GET",
			Description: "See a block",
			Payload:     "Block hash",
		},
	}
	json.NewEncoder(rw).Encode(data)

}
