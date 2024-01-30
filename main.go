package main

import (
	"fmt"
	"log"

	"github.com/jagottsicher/myGoBlockchain/wallet"
)

func init() {
	log.SetPrefix("Blockchain Node: ")
}

func main() {
	w := wallet.NewWallet()
	fmt.Println(w.PrivateKeyStr())
	fmt.Println(w.PublicKeyStr())
	fmt.Println(w.BlockchainAddress())

	t := wallet.NewTransaction(w.PrivateKey(), w.PublicKey(), w.BlockchainAddress(), "Morty", 137.0)
	fmt.Printf("Signature %s", t.GenerateSignature())

}
