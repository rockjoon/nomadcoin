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

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

type addTxPayload struct {
	To     string
	Amount int
}

type balanceResponse struct {
	Address string `json:"owner"`
	Balance int    `json:"balance"`
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
	router.HandleFunc("/status", status).Methods(http.MethodGet)
	router.HandleFunc("/blocks", blocks).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods(http.MethodGet)
	router.HandleFunc("/balance/{address}", balance).Methods(http.MethodGet)
	router.HandleFunc("/mempool", mempool).Methods(http.MethodGet)
	router.HandleFunc("/transactions", transactions).Methods(http.MethodGet, http.MethodPost)
	log.Printf("listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func transactions(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		break
	case http.MethodPost:
		var payload addTxPayload
		utils.HandleError(json.NewDecoder(r.Body).Decode(&payload))
		err := blockchain.Mempool.AddTxs(payload.To, payload.Amount)
		if err != nil {
			json.NewEncoder(rw).Encode(errorResponse{fmt.Sprint(err)})
			rw.WriteHeader(http.StatusBadRequest)
		} else {
			rw.WriteHeader(http.StatusCreated)
		}
	}
}

func mempool(rw http.ResponseWriter, r *http.Request) {
	json.NewEncoder(rw).Encode(blockchain.Mempool.Txs)
}

func balance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	total := r.URL.Query().Get("total")
	switch total {
	case "true":
		balance := blockchain.BalanceByAddress(address, blockchain.GetBlockChain())
		json.NewEncoder(rw).Encode(balanceResponse{address, balance})
		break
	default:
		json.NewEncoder(rw).Encode(blockchain.UTxOutsByAddress(address, blockchain.GetBlockChain()))
	}
}

func status(rw http.ResponseWriter, r *http.Request) {
	json.NewEncoder(rw).Encode(blockchain.GetBlockChain())
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
		json.NewEncoder(rw).Encode(blockchain.AllBlocks(blockchain.GetBlockChain()))
	case http.MethodPost:
		blockchain.GetBlockChain().AddBlock()
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
			Url:         URL("/status"),
			Method:      "GET",
			Description: "See the Blockchain status",
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
		{
			Url:         URL("/balance/{address}"),
			Method:      "GET",
			Description: "See a balance of address",
			Payload:     "address",
		},
	}
	json.NewEncoder(rw).Encode(data)

}
