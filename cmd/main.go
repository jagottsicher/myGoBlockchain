package main

import (
	"fmt"

	"github.com/jagottsicher/myGoBlockchain/utils"
)

func main() {
	// fmt.Println(utils.IsFoundNode("127.0.0.1", 3333))
	// fmt.Println(utils.IsFoundNode("localhost", 3333))

	fmt.Println(utils.FindNeighbors("127.0.0.1", 3333, 0, 3, 3333, 3336))
}
