package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/jagottsicher/myGoBlockchain/blockchain"
	block "github.com/jagottsicher/myGoBlockchain/blockchain"
	"github.com/jagottsicher/myGoBlockchain/utils"
	"github.com/jagottsicher/myGoBlockchain/wallet"
)

var cache map[string]*blockchain.Blockchain = make(map[string]*blockchain.Blockchain)

type BlockchainNode struct {
	port uint16
}

func NewBlockchainNode(port uint16) *BlockchainNode {
	return &BlockchainNode{port}
}

func (bcn *BlockchainNode) Port() uint16 {
	return bcn.port
}

func (bcn *BlockchainNode) GetBlockchain() *blockchain.Blockchain {
	bc, ok := cache["blockchain"]
	if !ok {
		minerWallet := wallet.NewWallet()
		bc = blockchain.NewBlockchain(minerWallet.BlockchainAddress(), bcn.Port())
		cache["blockchain"] = bc
		log.Printf("Public Key %v", minerWallet.PublicKeyStr())
		log.Printf("Private Key %v", minerWallet.PrivateKeyStr())
		log.Printf("Blockchain Address %v", minerWallet.BlockchainAddress())
	}
	return bc
}

func (bcn *BlockchainNode) GetChain(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		bc := bcn.GetBlockchain()
		m, _ := bc.MarshalJSON()
		io.WriteString(w, string(m[:]))
	default:
		log.Printf("Error: Invalid HTTP method")
	}

}

func (bcn *BlockchainNode) Transactions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		bc := bcn.GetBlockchain()
		transactions := bc.TransactionPool()
		m, _ := json.Marshal(struct {
			Transactions []*block.Transaction `json:"transactions"`
			Length       int                  `json:"length"`
		}{
			Transactions: transactions,
			Length:       len(transactions),
		})
		io.WriteString(w, string(m[:]))

	case http.MethodPost:
		decoder := json.NewDecoder((r.Body))
		var t blockchain.TransactionRequest
		err := decoder.Decode(&t)
		if err != nil {
			log.Printf("ERROR: %v", err)
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}
		if !t.Validate() {
			log.Println("ERROR: missing field(S)")
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}
		publicKey := utils.String2PublicKey(*t.SenderPublicKey)
		signature := utils.String2Signature(*t.Signature)
		bc := bcn.GetBlockchain()
		isCreated := bc.CreateTransaction(*t.SenderBlockchainAddress, *t.RecipientBlockchainAddress, *t.Value, publicKey, signature)

		w.Header().Add("Content-Type", "application/json")
		var m []byte
		if !isCreated {
			w.WriteHeader(http.StatusBadRequest)
			m = utils.JsonStatus("fail")
		} else {
			w.WriteHeader(http.StatusCreated)
			m = utils.JsonStatus("success")
		}
		io.WriteString(w, string(m))

	default:
		log.Println("ERROR: Invalid HTTP method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bcn *BlockchainNode) Mine(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		bc := bcn.GetBlockchain()
		isMined := bc.Mining()

		var m []byte
		if !isMined {
			w.WriteHeader(http.StatusBadRequest)
			m = utils.JsonStatus("fail")
		} else {
			m = utils.JsonStatus("success")
		}
		w.Header().Add("Content-Type", "application/json")
		io.WriteString(w, string(m))
	default:
		log.Println("ERROR: Invalid HTTP method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bcn *BlockchainNode) Run() {
	http.HandleFunc("/", bcn.GetChain)
	http.HandleFunc("/transactions", bcn.Transactions)
	http.HandleFunc("/mine", bcn.Mine)

	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(bcn.Port())), nil))
}
